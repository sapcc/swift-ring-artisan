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
	"slices"
	"sort"

	"github.com/sapcc/go-bits/logg"

	"github.com/sapcc/swift-ring-artisan/pkg/builderfile"
)

// NodeRules is a server containing disks
type NodeRules struct {
	Port           uint64             `yaml:"port,omitempty"`
	Meta           *map[string]string `yaml:"meta,omitempty"`
	DiskCount      uint64             `yaml:"disk_count"`
	DiskSizeTB     float64            `yaml:"disk_size_tb,omitempty"`
	Weight         *float64           `yaml:"weight,omitempty"`
	ReportedWeight *float64           `yaml:"reported_weight,omitempty"`
	// BrokenDisks lists device names like "swift-02" that shall be treated as non-existent.
	BrokenDisks []string `yaml:"broken_disks,omitempty"`
}

func (nodeRules NodeRules) DesiredWeight(baseSizeTB float64, nodeIP string) float64 {
	var weight float64
	switch {
	case nodeRules.Weight == nil && baseSizeTB == 0:
		logg.Fatal("Applying rule %+v to disk %s:%d failed because not enough data is present to calculate the weight", nodeRules, nodeIP, nodeRules.Port)
	case nodeRules.Weight == nil && baseSizeTB != 0:
		if nodeRules.DiskSizeTB == 0 {
			weight = 100
		} else {
			weight = math.Floor(nodeRules.DiskSizeTB / baseSizeTB * 100)
		}
	default:
		weight = *nodeRules.Weight
	}

	if weight == 0 {
		logg.Info("node.Weight %+v ruleData.BaseSizeTB %+v", nodeRules.Weight, baseSizeTB)
	}

	return weight
}

// ZoneRules contains multiple nodes
type ZoneRules struct {
	Nodes map[string]*NodeRules
}

func (zoneRules ZoneRules) getNodeIPs() []string {
	var nodeIPs []string
	for nodeIP := range zoneRules.Nodes {
		nodeIPs = append(nodeIPs, nodeIP)
	}
	sort.Strings(nodeIPs) // for reproducibility in tests

	return nodeIPs
}

type discoveredDisk struct {
	NodeIP   string
	DiskPort uint64
	DiskName string
}

func getDiscoveredDisk(nodeIP string, diskPort uint64, diskName string) discoveredDisk {
	return discoveredDisk{NodeIP: nodeIP, DiskPort: diskPort, DiskName: diskName}
}

// RingRules containing the rules for a region, multiple Zones and dozzens Nodes
type RingRules struct {
	BaseSizeTB float64 `yaml:"base_size_tb"`
	BasePort   uint64  `yaml:"base_port"`
	Region     uint64
	Overload   float64
	Zones      map[uint64]*ZoneRules
}

func (ringRules RingRules) getZones() []uint64 {
	var zones []uint64

	for zone := range ringRules.Zones {
		zones = append(zones, zone)
	}
	slices.Sort(zones)

	return zones
}

// CalculateChanges to parsed MetaData
func (ringRules RingRules) CalculateChanges(ring builderfile.RingInfo, ringFilename string) (commandQueue, confirmations []string, err error) {
	if ring.Regions == 0 {
		return nil, nil, errors.New("regions needs to be set")
	} else if ringRules.Region != ring.Regions || ring.Regions != 1 {
		return nil, nil, errors.New("currently only one region is supported")
	}

	var discoveredDisks []discoveredDisk

	// Special handling for floating point comparison
	if diff := math.Abs(ring.OverloadFactorDecimal - ringRules.Overload); diff > 0.000001 {
		logg.Debug("Overload does not match, adding command to change it from %f to %f", ring.OverloadFactorDecimal, ringRules.Overload)
		commandQueue = append(commandQueue, ring.CommandSetOverload(ringFilename, ringRules.Overload))
	}

	zones := ringRules.getZones()
	for _, zone := range zones {
		zoneRules := ringRules.Zones[zone]

		for _, nodeIP := range zoneRules.getNodeIPs() {
			nodeRules := zoneRules.Nodes[nodeIP]

			for diskNumber := uint64(1); diskNumber <= nodeRules.DiskCount; diskNumber++ {
				diskName := fmt.Sprintf("swift-%02d", diskNumber)
				if slices.Contains(nodeRules.BrokenDisks, diskName) {
					continue
				}

				weight := nodeRules.DesiredWeight(ringRules.BaseSizeTB, nodeIP)
				var port uint64
				switch {
				case nodeRules.Port != 0:
					port = nodeRules.Port
				case ringRules.BasePort != 0:
					port = ringRules.BasePort
				default:
					port = 6000
				}
				disk, err := ring.FindDevice(zone, nodeIP, port, diskName)
				if err != nil {
					return nil, nil, err
				}

				if disk == nil {
					logg.Debug("Disk was not found, adding it")
					disk = &builderfile.DeviceInfo{
						Region: ringRules.Region,
						Zone:   zone,
						NodeIP: nodeIP,
						Port:   port,
						Name:   diskName,
						Weight: weight,
					}
					if nodeRules.Meta != nil {
						disk.Meta = nodeRules.Meta
					}
					commandQueue = append(commandQueue, disk.CommandAdd(ringFilename))
					continue
				}

				discoveredDisks = append(discoveredDisks, getDiscoveredDisk(nodeIP, disk.Port, disk.Name))

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
		if !slices.Contains(discoveredDisks, getDiscoveredDisk(device.NodeIP, device.Port, device.Name)) {
			isNodeRemoved := false
			for _, zone := range zones {
				nodeIPs := ringRules.Zones[zone].getNodeIPs()
				if !slices.Contains(nodeIPs, device.NodeIP) {
					isNodeRemoved = true
				}
			}

			if device.Weight != 0 && isNodeRemoved {
				msg := fmt.Sprintf("Do you want to remove disk %s on node %s without first scaling its weight to 0? This poses a data loss risk.", device.Name, device.NodeIP)
				confirmations = append(confirmations, msg)
			}

			commandQueue = append(commandQueue, device.CommandRemove(ringFilename))
		}
	}

	return commandQueue, confirmations, nil
}
