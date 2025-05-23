// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package builderfile

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/oriser/regroup"
	"github.com/sapcc/go-bits/logg"

	"github.com/sapcc/swift-ring-artisan/pkg/misc"
)

// regex to match the following line:
// container.builder, build version 7, id 024e79c994c643d09eb045d488dafb94
// account.builder, build version 37, id (not assigned)
var fileInfoRx = regroup.MustCompile(`^(?:[\w\/\.-]+\/)?(?P<fileName>\w+\.builder), build version (?P<buildVersion>\d+), id (?:(?P<id>[\d\w]{32})|\(not assigned\))$`)

// regex to match the following line:
// 1024 partitions, 3.000000 replicas, 1 regions, 1 zones, 6 devices, 0.00 balance, 0.00 dispersion
var statsRx = regroup.MustCompile(`^(?P<partitions>\d+) partitions, (?P<replicas>\d+\.\d+) replicas, (?P<regions>\d+) regions, (?P<zones>\d) zones, (?P<deviceCount>\d+) devices, (?P<balance>\d+\.\d+) balance, (?P<dispersion>\d+\.\d+) dispersion$`)

// regex to match the following line:
// The minimum number of hours before a partition can be reassigned is 24 (0:00:00 remaining)
var remainingTimeRx = regroup.MustCompile(`^The minimum number of hours before a partition can be reassigned is (?P<reassignedCooldown>\d+) \((?P<reassignedRemaining>\d+:\d{2}:\d{2}) remaining\)$`)

// regex to match the following line:
// The overload factor is 0.00% (0.000000)
var overloadFactorRx = regroup.MustCompile(`^The overload factor is (?P<percent>\d+\.\d+)% \((?P<decimal>\d+\.\d+)\)$`)

// regex to match the following line:
// Ring file container.ring.gz is obsolete
// Ring file container.ring.gz is up-to-date
var obsoleteRx = regexp.MustCompile(`^Ring file (?:[\w\/\.-]+\/)?\w+\.ring\.gz is (obsolete|up-to-date)$`)

// regex to match the following line:
// Devices:   id region zone   ip address:port replication ip:port  name weight partitions balance flags meta
var tableHeaderRx = regexp.MustCompile(`^Devices:   id region zone\s+ip address:port replication ip:port  name weight partitions balance flags meta$`)

// regex to match the following line:
//
//	  0      1    1 10.114.1.202:6001   10.114.1.202:6001 swift-01 100.00        512    0.00
//	  1      1    1 10.114.1.202:6001   10.114.1.202:6001 swift-02 100.00        512    0.00
//	  2      1    1 10.114.1.202:6001   10.114.1.202:6001 swift-03 100.00        512    0.00
//	111      1    1  10.46.14.44:6001    10.46.14.44:6001 swift-33 100.00         78   -0.98
//	 65      1    1   10.46.14.44:6002    10.46.14.44:6002 swift-01 100.00         64   -5.63       {"hostname":"nodeswift01-cp001"}
var rowEntryRx = regroup.MustCompile(`^\s+(?P<id>\d+)\s+(?P<region>\d+)\s+(?P<zone>\d+)\s+(?P<ip>(?:\d+\.){3}\d+):(?P<port>\d+)\s+(?P<replicationIp>(?:\d+\.){3}\d+):(?P<replicationPort>\d+)\s+(?P<name>[\w+-]+)\s+(?P<weight>\d+\.\d+)\s+(?P<partitions>\d+)\s+(?P<balance>-?\d+\.\d+)\s*(?P<meta>\{"hostname":"\w+-\w+"\})?$`)

// FindDevice returns a given disk that matches the in
func (ring RingInfo) FindDevice(zone uint64, nodeIP string, port uint64, diskName string) (*DeviceInfo, error) {
	for _, dev := range ring.Devices {
		// zone is not checked here to detect potential zone mismatches
		// if there are ever nodes which split disks across multiple zones this will break
		// if zone would be checked here a command to remove and add a disk on a zone mismatch would be generated
		if dev.NodeIP == nodeIP && dev.Name == diskName {
			if dev.Zone != zone {
				return nil, fmt.Errorf("zone ID mismatch between parsed data %d and rule file %d", dev.Zone, zone)
			}
			if dev.Port != port {
				return nil, fmt.Errorf("port mismatch between parsed data %d and rule file %d", dev.Port, port)
			}
			if dev.NodeIP != dev.ReplicationIP || dev.Port != dev.ReplicationPort {
				return nil, errors.New("replication ip and port do not match with the normal ip and port which is required")
			}
			return &dev, nil
		}
	}

	return nil, nil
}

// Input parses an input and return the data as MetData object
func Input(input io.Reader) RingInfo {
	var metaData RingInfo
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		logg.Debug("Processing line: %s\n", line)

		//nolint:errcheck // errors can be ignored because the regex matches digits (\d)
		matches, _ := fileInfoRx.Groups(line)
		if len(matches) > 0 {
			metaData.FileName = matches["fileName"]
			metaData.Version = misc.ParseUint(matches["buildVersion"])
			metaData.ID = matches["id"]
			continue
		}

		//nolint:errcheck // errors can be ignored because the regex matches digits (\d)
		matches, _ = statsRx.Groups(line)
		if len(matches) > 0 {
			metaData.Partitions = misc.ParseUint(matches["partitions"])
			metaData.Replicas = misc.ParseFloat(matches["replicas"])
			metaData.Regions = misc.ParseUint(matches["regions"])
			metaData.Zones = misc.ParseUint(matches["zones"])
			metaData.DeviceCount = misc.ParseUint(matches["deviceCount"])
			metaData.Balance = misc.ParseFloat(matches["balance"])
			metaData.Dispersion = misc.ParseFloat(matches["dispersion"])
			continue
		}

		//nolint:errcheck // errors can be ignored because the regex matches digits (\d)
		matches, _ = remainingTimeRx.Groups(line)
		if len(matches) > 0 {
			metaData.ReassignedCooldown = misc.ParseUint(matches["reassignedCooldown"])
			metaData.ReassignedRemaining, _ = time.Parse(time.TimeOnly, matches["reassignedRemaining"]) //nolint:errcheck
			continue
		}

		matches, _ = overloadFactorRx.Groups(line) //nolint:errcheck
		if len(matches) > 0 {
			// errors can be ignored because the regex matches digits (\d)
			metaData.OverloadFactorPercent = misc.ParseFloat(matches["percent"])
			metaData.OverloadFactorDecimal = misc.ParseFloat(matches["decimal"])
			continue
		}

		// this line is purely informational but we need to match it anyway to not abort the process
		if obsoleteRx.MatchString(line) {
			continue
		}

		if tableHeaderRx.MatchString(line) {
			break
		}

		logg.Fatal("A header regex did not match the line: \"%s\"", line)
	}

	for scanner.Scan() {
		line := scanner.Text()
		logg.Debug("Processing line: %s\n", line)

		matches, _ := rowEntryRx.Groups(line) //nolint:errcheck
		if len(matches) > 0 {
			// errors can be ignored because the regex matches digits (\d)
			var meta *map[string]string
			if matches["meta"] != "" {
				// logg.Info("%#v", matches["meta"])
				err := json.Unmarshal([]byte(matches["meta"]), &meta)
				if err != nil {
					logg.Fatal("Unmsharshaling meta from swift-ring-builder output failed: %s %s", matches["meta"], err.Error())
				}
			}

			metaData.Devices = append(metaData.Devices, DeviceInfo{
				ID:              misc.ParseUint(matches["id"]),
				Region:          misc.ParseUint(matches["region"]),
				Zone:            misc.ParseUint(matches["zone"]),
				NodeIP:          matches["ip"],
				Port:            misc.ParseUint(matches["port"]),
				ReplicationIP:   matches["replicationIp"],
				ReplicationPort: misc.ParseUint(matches["replicationPort"]),
				Name:            matches["name"],
				Weight:          misc.ParseFloat(matches["weight"]),
				Partitions:      misc.ParseUint(matches["partitions"]),
				// disabled because the information cannot easily be extracted from the pickle file
				// which causes mismatches when comparing the outputs
				// Balance:         misc.ParseFloat(matches["balance"]),
				Meta: meta,
			})
			continue
		}

		logg.Fatal("The table entry regex did not match the line: %s", line)
	}

	if err := scanner.Err(); err != nil {
		logg.Fatal("Reading input failed: %s", err.Error())
	}

	return metaData
}
