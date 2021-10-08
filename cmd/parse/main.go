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

package parsecmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/oriser/regroup"
	"github.com/sapcc/go-bits/logg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	format string
)

func AddCommandTo(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "parse <image>",
		Example: "  swift-ring-builder object.builder | swift-ring-artisan parse",
		Short:   "Parses the output of swift-ring-builder to a machine readable format.",
		Long: `Parses the output of swift-ring-builder to a machine readable format.
The output of swift-ring-builder needs to be piped into swift-ring-artisan.`,
		Run: run,
	}
	cmd.PersistentFlags().StringVar(&format, "format", "[]", "Output format. Can be either json or yaml.")
	parent.AddCommand(cmd)
}

// regex to match the following line:
// container.builder, build version 7, id 024e79c994c643d09eb045d488dafb94
var fileInfoRx = regroup.MustCompile(`^(?P<fileName>\w+\.builder), build version (?P<buildVersion>\d+), id (?P<id>[\d\w]{32})$`)

// regex to match the following line:
// 1024 partitions, 3.000000 replicas, 1 regions, 1 zones, 6 devices, 0.00 balance, 0.00 dispersion
var statsRx = regroup.MustCompile(`^(?P<partitions>\d+) partitions, (?P<repliacs>\d+\.\d+) replicas, (?P<regions>\d+) regions, (?P<zones>\d) zones, (?P<deviceCount>\d+) devices, (?P<balance>\d+\.\d+) balance, (?P<dispersion>\d+\.\d+) dispersion$`)

// regex to match the following line:
// The minimum number of hours before a partition can be reassigned is 24 (0:00:00 remaining)
var remainingTimeRx = regroup.MustCompile(`^The minimum number of hours before a partition can be reassigned is (?P<reassignedCooldown>\d+) \((?P<reassignedRemaining>\d+:\d{2}:\d{2}) remaining\)$`)

// regex to match the following line:
// The overload factor is 0.00% (0.000000)
var overloadFactorRx = regroup.MustCompile(`^The overload factor is (?P<percent>\d+\.\d+)% \((?P<dezimal>\d+\.\d+\))$`)

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
	Id                int
	Region            int
	Zone              int
	IpAddressPort     string
	ReplicationIpPort string
	Name              string
	Weight            float64
	Partitions        int
	Balance           int
	flags             struct{} // TODO: figure out how the field looks like
	meta              struct{} // TODO: figure out how the field looks like
}

// Meta data about the ring file
type MetaData struct {
	FileName     string
	BuildVersion int
	Id           string

	Partitions  int
	Replicas    float64
	Regions     int
	Zones       int
	DeviceCount int
	Balance     float64
	Dispersion  float64

	ReassignedCooldown  int
	ReassignedRemaining time.Time

	OverloadFactorPercent float64
	OverloadFactorDezimal float64

	Devices []device
}

func run(cmd *cobra.Command, args []string) {
	// if a file is supplied read that otherwise listen on stdin
	var input io.Reader
	if len(args) == 1 {
		file, err := os.Open(args[0])
		if err != nil {
			logg.Fatal(fmt.Sprintf("Reading file failed: %s", err.Error()))
		}
		defer file.Close()
		input = file
	} else {
		input = os.Stdin
	}

	var (
		metaData MetaData
		// track if we processed all headers and then only match table entries
		processedTableEntries bool
	)
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		logg.Debug(fmt.Sprintf("Processing line: %s\n", line))

		if !processedTableEntries {
			matches, _ := fileInfoRx.Groups(line)
			if len(matches) > 0 {
				metaData.FileName = matches["fileName"]
				metaData.BuildVersion, _ = strconv.Atoi(matches["buildVersion"])
				metaData.Id = matches["id"]
				continue
			}

			matches, _ = remainingTimeRx.Groups(line)
			if len(matches) > 0 {
				metaData.Partitions, _ = strconv.Atoi(matches["partitions"])
				metaData.Replicas, _ = strconv.ParseFloat(matches["replicas"], 32)
				metaData.Regions, _ = strconv.Atoi(matches["regions"])
				metaData.Zones, _ = strconv.Atoi(matches["zones"])
				metaData.DeviceCount, _ = strconv.Atoi(matches["deviceCount"])
				metaData.Balance, _ = strconv.ParseFloat(matches["balance"], 32)
				metaData.Dispersion, _ = strconv.ParseFloat(matches["dispersion"], 32)
				continue
			}

			matches, _ = statsRx.Groups(line)
			if len(matches) > 0 {
				metaData.ReassignedCooldown, _ = strconv.Atoi(matches["reassignedCooldown"])
				metaData.ReassignedRemaining, _ = time.Parse("15:04:05", matches["reassignedRemaining"])
				continue
			}

			matches, _ = overloadFactorRx.Groups(line)
			if len(matches) > 0 {
				metaData.OverloadFactorPercent, _ = strconv.ParseFloat(matches["percent"], 32)
				metaData.OverloadFactorDezimal, _ = strconv.ParseFloat(matches["dezimal"], 32)
				continue
			}

			// this line is purely informationaly but we need to match it anyway to not abort the process
			if obsoleteRx.MatchString(line) {
				continue
			}

			if tableHeaderRx.MatchString(line) {
				processedTableEntries = true
				continue
			}

			fmt.Printf("%+v\n", metaData)

			logg.Fatal(fmt.Sprintf("A header regex did not match the line: \"%s\"", line))
		} else {
			matches, _ := rowEntryRx.Groups(line)
			if len(matches) > 0 {
				id, _ := strconv.Atoi(matches["id"])
				region, _ := strconv.Atoi(matches["region"])
				zone, _ := strconv.Atoi(matches["zone"])
				weight, _ := strconv.ParseFloat(matches["weight"], 32)
				partitions, _ := strconv.Atoi(matches["partitions"])
				balance, _ := strconv.Atoi(matches["balance"])

				metaData.Devices = append(metaData.Devices, device{
					Id:                id,
					Region:            region,
					Zone:              zone,
					IpAddressPort:     matches["ipAddressPort"],
					ReplicationIpPort: matches["replicationIpPort"],
					Name:              matches["name"],
					Weight:            weight,
					Partitions:        partitions,
					Balance:           balance,
				})
				continue
			}

			logg.Fatal(fmt.Sprintf("The table entry regex did not match the line: %s", line))
		}
	}

	err := scanner.Err()
	if err != nil {
		logg.Fatal(fmt.Sprintf("Reading input failed: %s", err.Error()))
	}

	metaDataYaml, err := yaml.Marshal(metaData)
	if err != nil {
		logg.Fatal(err.Error())
	}

	fmt.Printf("%+v", string(metaDataYaml))

	if err != nil {
		os.Exit(1)
	}
}
