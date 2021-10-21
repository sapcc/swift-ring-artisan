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

package rules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/oriser/regroup"
	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/swift-ring-artisan/pkg/parse"
)

// Disk properties
type Disk struct {
	Count  uint64
	Size   string
	Weight *float64 `yaml:"weight,omitempty"`
}

// Node is a server containing disks
type Node struct {
	IPPort string   `yaml:"ip_port"`
	Meta   struct{} `yaml:"meta,omitempty"` // TODO: figure out how the field looks like
	Disks  Disk
}

// Zone contains multiple nodes
type Zone struct {
	ID    uint64
	Nodes []Node
}

// DiskRules containing the rules for a region, multiple Zones and dozzens Nodes
type DiskRules struct {
	BaseSize string `yaml:"base_size"`
	Region   uint64
	Zones    []Zone
}

var storageRx = regroup.MustCompile(`(?P<size>\d+)(?P<unit>\w+)`)

// ApplyRules to parsed MetaData
func ApplyRules(inputData parse.MetaData, ruleData DiskRules, ringFilename string) []string {
	if ruleData.Region != inputData.Regions || inputData.Regions != 1 {
		logg.Fatal("Only one region is currently supported.")
	}

	var commandQueue []string
	counter := 0
	for i := range ruleData.Zones {
		zone := ruleData.Zones[i]
		if zone.ID != uint64(i+1) {
			logg.Fatal("Zone ID mismatch between parsed data and rule file.")
		}

		for j := range zone.Nodes {
			node := zone.Nodes[j]

			for k := 0; k < int(node.Disks.Count); k++ {
				diskRules := node.Disks
				diskData := inputData.Devices[counter]
				logg.Debug(fmt.Sprintf("Applying rule %+v to disk number %d: %+v", diskRules, counter, diskData))

				if node.IPPort != diskData.IPAddressPort {
					logg.Fatal(fmt.Sprintf("The IP port combination of the rule number %d does not match the by id sorted parsed one: %s != %s",
						counter, node.IPPort, diskData.IPAddressPort))
				}

				if diskRules.Weight == nil && diskRules.Size != "" {
					matches, _ := storageRx.Groups(diskRules.Size)
					if len(matches) == 0 {
						logg.Fatal(fmt.Sprintf("Can't parse size into a value unit pair: %s", diskRules.Size))
					}
					size, _ := strconv.ParseFloat(matches["size"], 32)
					matches, _ = storageRx.Groups(ruleData.BaseSize)
					if len(matches) == 0 {
						logg.Fatal(fmt.Sprintf("Can't parse baseSize into a value unit pair: %s", ruleData.BaseSize))
					}
					baseSize, _ := strconv.ParseFloat(matches["size"], 32)
					// poor mans rounding to convert 166.666 to 166 instead of 167
					weight := float64(int(size / baseSize * 100))
					diskRules.Weight = &weight
				}

				if diskRules.Weight == nil {
					logg.Fatal(fmt.Sprintf("Applying rule %+v to disk number %d failed because weight is not set", diskRules, counter))
				}

				if *diskRules.Weight != diskData.Weight {
					logg.Debug("Weight does not match, adding command to change it")
					ipAddressPort := strings.Split(diskData.IPAddressPort, ":")
					commandQueue = append(commandQueue, fmt.Sprintf(
						"swift-ring-builder %s set_weight --region %d --zone %d --ip %s --port %s --device %s --weight %g %g",
						ringFilename, diskData.Region, diskData.Zone, ipAddressPort[0], ipAddressPort[1], diskData.Name, diskData.Weight, *diskRules.Weight))
				}
				counter++
			}
		}
	}

	return commandQueue
}
