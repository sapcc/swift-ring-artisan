/*******************************************************************************
*
* Copyright 2021 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

package parse

import (
	"bufio"
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
var fileInfoRx = regroup.MustCompile(`^(?P<fileName>\w+\.builder), build version (?P<buildVersion>\d+), id (?P<id>[\d\w]{32})$`)

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
var obsoleteRx = regexp.MustCompile(`^Ring file \w+\.ring\.gz is obsolete$`)

// regex to match the following line:
// Devices:   id region zone   ip address:port replication ip:port  name weight partitions balance flags meta
var tableHeaderRx = regexp.MustCompile(`^Devices:   id region zone   ip address:port replication ip:port  name weight partitions balance flags meta$`)

// regex to match the following line:
//            0      1    1 10.114.1.202:6001   10.114.1.202:6001 swift-01 100.00        512    0.00
//            1      1    1 10.114.1.202:6001   10.114.1.202:6001 swift-02 100.00        512    0.00
//            2      1    1 10.114.1.202:6001   10.114.1.202:6001 swift-03 100.00        512    0.00
var rowEntryRx = regroup.MustCompile(`^\s+(?P<id>\d+)\s+(?P<region>\d+)\s+(?P<zone>\d+)\s+(?P<ipAddressPort>(?:\d+\.){3}\d+:\d+)\s+(?P<replicationIpPort>(?:\d+\.){3}\d+:\d+)\s+(?P<name>[\w+-]+)\s+(?P<weight>\d+\.\d+)\s+(?P<partitions>\d+)\s+(?P<balance>\d+\.\d+)\s*$`)

type device struct {
	ID                uint64
	Region            uint64
	Zone              uint64
	IPAddressPort     string `yaml:"ip_address_port"`
	ReplicationIPPort string `yaml:"replication_ip_port"`
	Name              string
	Weight            float64
	Partitions        uint64
	Balance           float64
	//lint:ignore U1000 TODO
	flags struct{} // TODO: figure out how the field looks like
	//lint:ignore U1000 TODO
	meta struct{} // TODO: figure out how the field looks like
}

// MetaData contains the meta data about the ring file
type MetaData struct {
	FileName     string `yaml:"file_name"`
	BuildVersion uint64 `yaml:"build_version"`
	ID           string

	Partitions  uint64
	Replicas    float64
	Regions     uint64
	Zones       uint64
	DeviceCount uint64 `yaml:"device_count"`
	Balance     float64
	Dispersion  float64

	ReassignedCooldown  uint64    `yaml:"reassigned_cooldown"`
	ReassignedRemaining time.Time `yaml:"reassigned_remaining"`

	OverloadFactorPercent float64 `yaml:"overload_factor_Percent"`
	OverloadFactorDecimal float64 `yaml:"overload_factor_decimal"`

	Devices []device
}

// Input parses an input and return the data as MetData object
func Input(input io.Reader) MetaData {
	var metaData MetaData

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		logg.Debug(fmt.Sprintf("Processing line: %s\n", line))

		matches, _ := fileInfoRx.Groups(line)
		if len(matches) > 0 {
			metaData.FileName = matches["fileName"]
			// errors can be ignored because the regex matches digits (\d)
			metaData.BuildVersion = misc.ParseUint(matches["buildVersion"])
			metaData.ID = matches["id"]
			continue
		}

		matches, _ = statsRx.Groups(line)
		if len(matches) > 0 {
			// errors can be ignored because the regex matches digits (\d)
			metaData.Partitions = misc.ParseUint(matches["partitions"])
			metaData.Replicas = misc.ParseFloat(matches["replicas"])
			metaData.Regions = misc.ParseUint(matches["regions"])
			metaData.Zones = misc.ParseUint(matches["zones"])
			metaData.DeviceCount = misc.ParseUint(matches["deviceCount"])
			metaData.Balance = misc.ParseFloat(matches["balance"])
			metaData.Dispersion = misc.ParseFloat(matches["dispersion"])
			continue
		}

		matches, _ = remainingTimeRx.Groups(line)
		if len(matches) > 0 {
			// errors can be ignored because the regex matches digits (\d)
			metaData.ReassignedCooldown = misc.ParseUint(matches["reassignedCooldown"])
			metaData.ReassignedRemaining, _ = time.Parse("15:04:05", matches["reassignedRemaining"])
			continue
		}

		matches, _ = overloadFactorRx.Groups(line)
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

		logg.Fatal(fmt.Sprintf("A header regex did not match the line: \"%s\"", line))
	}

	for scanner.Scan() {
		line := scanner.Text()
		logg.Debug(fmt.Sprintf("Processing line: %s\n", line))

		matches, _ := rowEntryRx.Groups(line)
		if len(matches) > 0 {
			// errors can be ignored because the regex matches digits (\d)
			metaData.Devices = append(metaData.Devices, device{
				ID:                misc.ParseUint(matches["id"]),
				Region:            misc.ParseUint(matches["region"]),
				Zone:              misc.ParseUint(matches["zone"]),
				IPAddressPort:     matches["ipAddressPort"],
				ReplicationIPPort: matches["replicationIpPort"],
				Name:              matches["name"],
				Weight:            misc.ParseFloat(matches["weight"]),
				Partitions:        misc.ParseUint(matches["partitions"]),
				Balance:           misc.ParseFloat(matches["balance"]),
			})
			continue
		}

		logg.Fatal(fmt.Sprintf("The table entry regex did not match the line: %s", line))
	}

	err := scanner.Err()
	if err != nil {
		logg.Fatal(fmt.Sprintf("Reading input failed: %s", err.Error()))
	}

	return metaData
}
