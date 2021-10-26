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
	"math"

	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/swift-ring-artisan/pkg/parse"
)

// Node is a server containing disks
type Node struct {
	Port       uint64   `yaml:"port,omitempty"`
	Meta       struct{} `yaml:"meta,omitempty"` // TODO: figure out how the field looks like
	DiskCount  uint64   `yaml:"disk_count"`
	DiskSizeTB float64  `yaml:"disk_size_tb,omitempty"`
	Weight     *float64 `yaml:"weight,omitempty"`
}

// Zone contains multiple nodes
type Zone struct {
	Zone  uint64
	Nodes map[string]Node
}

// DiskRules containing the rules for a region, multiple Zones and dozzens Nodes
type DiskRules struct {
	BaseSizeTB float64 `yaml:"base_size_tb"`
	BasePort   uint64  `yaml:"base_port"`
	Region     uint64
	Zones      []Zone
}

// ApplyRules to parsed MetaData
func ApplyRules(inputData parse.MetaData, ruleData DiskRules, ringFilename string) []string {
	if inputData.Regions == 0 {
		logg.Fatal("Regions needs to be set.")
	} else if ruleData.Region != inputData.Regions || inputData.Regions != 1 {
		logg.Fatal("Only one region is currently supported.")
	}

	var commandQueue []string
	for i, zone := range ruleData.Zones {
		if zone.Zone != uint64(i+1) {
			logg.Fatal("Zone ID mismatch between parsed data and rule file.")
		}

		for ip, node := range zone.Nodes {
			for diskNumber := 1; diskNumber <= int(node.DiskCount); diskNumber++ {
				var disk parse.Device
				diskName := fmt.Sprintf("swift-%02d", diskNumber)
				for _, dev := range inputData.Devices {
					if dev.IP == ip && dev.Name == diskName {
						disk = dev
						break
					}
				}
				if disk.Name == "" {
					logg.Fatal("No device found for ip %s and name %s", ip, diskName)
				}

				logg.Debug("Applying rule %+v to disk %s:%d %+v", node, ip, node.Port, disk)
				var weight float64
				if node.Weight == nil && ruleData.BaseSizeTB == 0 {
					logg.Fatal("Applying rule %+v to disk %s:%d failed because not enough data is present to calculate the weight", node, ip, node.Port)
				} else if node.Weight == nil && ruleData.BaseSizeTB != 0 {
					if node.DiskSizeTB == 0 {
						weight = 100
					} else {
						weight = math.Floor(node.DiskSizeTB / ruleData.BaseSizeTB * 100)
					}
				} else {
					weight = *node.Weight
				}

				if weight == 0 {
					logg.Info("node.Weight %+v ruleData.BaseSizeTB %+v", node.Weight, ruleData.BaseSizeTB)
				}

				if disk.Weight != weight {
					logg.Debug("Weight does not match, adding command to change it")
					commandQueue = append(commandQueue, fmt.Sprintf(
						"swift-ring-builder %s set_weight --region %d --zone %d --ip %s --port %d --device %s --weight %g %g",
						ringFilename, disk.Region, disk.Zone, disk.IP, disk.Port, disk.Name, disk.Weight, weight))
				}
			}
		}
	}

	return commandQueue
}
