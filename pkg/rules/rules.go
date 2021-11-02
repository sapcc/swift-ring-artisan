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
	"github.com/sapcc/swift-ring-artisan/pkg/builderfile"
)

// NodeRules is a server containing disks
type NodeRules struct {
	Port       uint64   `yaml:"port,omitempty"`
	Meta       struct{} `yaml:"meta,omitempty"` // TODO: figure out how the field looks like
	DiskCount  uint64   `yaml:"disk_count"`
	DiskSizeTB float64  `yaml:"disk_size_tb,omitempty"`
	Weight     *float64 `yaml:"weight,omitempty"`
}

// ZoneRules contains multiple nodes
type ZoneRules struct {
	Zone  uint64
	Nodes map[string]NodeRules
}

// RingRules containing the rules for a region, multiple Zones and dozzens Nodes
type RingRules struct {
	BaseSizeTB float64 `yaml:"base_size_tb"`
	BasePort   uint64  `yaml:"base_port"`
	Region     uint64
	Zones      []ZoneRules
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
func (ringRules RingRules) CalculateChanges(ring builderfile.RingInfo, ringFilename string) []string {
	if ring.Regions == 0 {
		logg.Fatal("Regions needs to be set.")
	} else if ringRules.Region != ring.Regions || ring.Regions != 1 {
		logg.Fatal("Only one region is currently supported.")
	}

	var commandQueue []string
	for zoneID, zoneRules := range ringRules.Zones {
		if zoneRules.Zone != uint64(zoneID+1) {
			logg.Fatal("Zone ID mismatch between parsed data and rule file.")
		}

		for nodeIP, nodeRules := range zoneRules.Nodes {
			for diskNumber := 1; diskNumber <= int(nodeRules.DiskCount); diskNumber++ {
				disk := ring.FindDevice(nodeIP, diskNumber)

				logg.Debug("Applying rule %+v to disk %s:%d %+v", nodeRules, nodeIP, nodeRules.Port, disk)
				weight := nodeRules.DesiredWeight(ringRules.BaseSizeTB, nodeIP)
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
