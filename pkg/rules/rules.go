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
	"errors"
	"fmt"
	"math"
	"reflect"
	"sort"

	"github.com/sapcc/go-bits/logg"

	"github.com/sapcc/swift-ring-artisan/pkg/builderfile"
	"github.com/sapcc/swift-ring-artisan/pkg/misc"
)

// NodeRules is a server containing disks
type NodeRules struct {
	Port           uint64             `yaml:"port,omitempty"`
	Meta           *map[string]string `yaml:"meta,omitempty"`
	DiskCount      uint64             `yaml:"disk_count"`
	DiskSizeTB     float64            `yaml:"disk_size_tb,omitempty"`
	Weight         *float64           `yaml:"weight,omitempty"`
	ReportedWeight *float64           `yaml:"reported_weight,omitempty"`
}

// ZoneRules contains multiple nodes
type ZoneRules struct {
	Nodes map[string]*NodeRules
}

// RingRules containing the rules for a region, multiple Zones and dozzens Nodes
type RingRules struct {
	BaseSizeTB float64 `yaml:"base_size_tb"`
	BasePort   uint64  `yaml:"base_port"`
	Region     uint64
	Overload   float64
	Zones      map[uint64]*ZoneRules
}

func (nodeRules NodeRules) DesiredWeight(baseSizeTB float64, nodeIP string) float64 {
	var weight float64
	if nodeRules.Weight == nil && baseSizeTB == 0 {
		logg.Fatal("Applying rule %+v to disk %s:%d failed because not enough data is present to calculate the weight", nodeRules, nodeIP, nodeRules.Port)
	} else if nodeRules.Weight == nil && baseSizeTB != 0 {
		if nodeRules.DiskSizeTB == 0 {
			weight = 100
		} else {
			weight = math.Floor(nodeRules.DiskSizeTB / baseSizeTB * 100)
		}
	} else {
		weight = *nodeRules.Weight
	}

	if weight == 0 {
		logg.Info("node.Weight %+v ruleData.BaseSizeTB %+v", nodeRules.Weight, baseSizeTB)
	}

	return weight
}

// CalculateChanges to parsed MetaData
func (ringRules RingRules) CalculateChanges(ring builderfile.RingInfo, ringFilename string) ([]string, error) {
	if ring.Regions == 0 {
		return nil, errors.New("regions needs to be set")
	} else if ringRules.Region != ring.Regions || ring.Regions != 1 {
		return nil, errors.New("currently only one region is supported")
	}

	var discoveredDisks, commandQueue []string

	if ring.OverloadFactorDecimal != ringRules.Overload {
		logg.Debug("Overload does not match, adding command to change it")
		commandQueue = append(commandQueue, ring.CommandSetOverload(ringFilename, ringRules.Overload))
	}

	var zoneIDs []uint64
	for zoneID := range ringRules.Zones {
		zoneIDs = append(zoneIDs, zoneID)
	}
	sort.Slice(zoneIDs, func(i, j int) bool { return zoneIDs[i] < zoneIDs[j] }) // for reproducibility in tests

	for _, zoneID := range zoneIDs {
		zoneRules := ringRules.Zones[zoneID]

		var nodeIPs []string
		for nodeIP := range zoneRules.Nodes {
			nodeIPs = append(nodeIPs, nodeIP)
		}
		sort.Strings(nodeIPs) // for reproducibility in tests

		for _, nodeIP := range nodeIPs {
			nodeRules := zoneRules.Nodes[nodeIP]
			for diskNumber := 1; diskNumber <= int(nodeRules.DiskCount); diskNumber++ {
				weight := nodeRules.DesiredWeight(ringRules.BaseSizeTB, nodeIP)
				var port uint64
				if nodeRules.Port != 0 {
					port = nodeRules.Port
				} else if ringRules.BasePort != 0 {
					port = ringRules.BasePort
				} else {
					port = 6000
				}
				disk, err := ring.FindDevice(zoneID, nodeIP, port, diskNumber)
				if err != nil {
					return nil, err
				}

				if disk == nil {
					logg.Debug("Disk was not found, adding it")
					disk = &builderfile.DeviceInfo{
						Region: ringRules.Region,
						Zone:   zoneID,
						IP:     nodeIP,
						Port:   port,
						Name:   fmt.Sprintf("swift-%02d", diskNumber),
						Weight: weight,
					}
					if nodeRules.Meta != nil {
						disk.Meta = nodeRules.Meta
					}
					commandQueue = append(commandQueue, disk.CommandAdd(ringFilename))
					continue
				}

				discoveredDisks = append(discoveredDisks, fmt.Sprintf("%s\000%d\000%s", nodeIP, disk.Port, disk.Name))

				logg.Debug("Applying rule %+v to disk %s:%d %+v", nodeRules, nodeIP, port, disk)
				if disk.Weight != weight {
					logg.Debug("Weight does not match, adding command to change it")
					commandQueue = append(commandQueue, disk.CommandSetWeight(ringFilename, weight))
				}

				if nodeRules.Meta != nil && !reflect.DeepEqual(disk.Meta, nodeRules.Meta) {
					logg.Debug("Meta does not match, adding command to change it")
					commandQueue = append(commandQueue, disk.CommandSetMeta(ringFilename, *nodeRules.Meta))
				}
			}
		}
	}

	// check if all devices in the ring where matched with a rule
	for _, device := range ring.Devices {
		if !misc.Contains(discoveredDisks, fmt.Sprintf("%s\000%d\000%s", device.IP, device.Port, device.Name)) {
			commandQueue = append(commandQueue, device.CommandRemove(ringFilename))
		}
	}

	return commandQueue, nil
}
