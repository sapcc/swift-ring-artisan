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

package convert

import (
	"github.com/sapcc/swift-ring-artisan/pkg/builderfile"
	"github.com/sapcc/swift-ring-artisan/pkg/rules"
)

// Convert converts parsed MetaData to DiskRules
func Convert(ring builderfile.RingInfo, basePort uint64, baseSize float64) rules.RingRules {
	diskRules := rules.RingRules{
		Region:     1, // FIXME: make multi region aware
		BasePort:   basePort,
		BaseSizeTB: baseSize,
	}

	var (
		last          string
		diskRulesZone *rules.ZoneRules
	)
	for _, device := range ring.Devices {
		// create zone if it does not exist
		if len(diskRules.Zones) == 0 || diskRules.Zones[len(diskRules.Zones)-1].Zone != device.Zone {
			diskRules.Zones = append(diskRules.Zones, rules.ZoneRules{Zone: device.Zone})
		}

		diskRulesZone = &diskRules.Zones[device.Zone-1]
		// if the last IPAddressPort matches the current, there is another disk on the same note, just increase the count
		if last != "" && last == device.IP {
			node := diskRulesZone.Nodes[device.IP]
			node.DiskCount++
			diskRulesZone.Nodes[device.IP] = node
			continue
		}
		if diskRulesZone.Nodes == nil {
			diskRulesZone.Nodes = make(map[string]rules.NodeRules)
		}
		diskRulesZone.Nodes[device.IP] = rules.NodeRules{
			DiskCount: 1,
			Weight:    &device.Weight,
		}

		last = device.IP
	}

	return diskRules
}
