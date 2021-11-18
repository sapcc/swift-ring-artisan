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
func Convert(ring builderfile.RingInfo, baseSize float64) rules.RingRules {
	diskRules := rules.RingRules{
		Region:     1, // FIXME: make multi region aware
		BasePort:   ring.Devices[0].Port,
		BaseSizeTB: baseSize,
		Zones:      make(map[uint64]*rules.ZoneRules),
	}

	var diskRulesZone *rules.ZoneRules
	for _, device := range ring.Devices {
		// create zone if it does not exist
		if _, ok := diskRules.Zones[device.Zone]; len(diskRules.Zones) == 0 || !ok {
			diskRules.Zones[device.Zone] = &rules.ZoneRules{}
		}

		diskRulesZone = diskRules.Zones[device.Zone]
		// if the last IPAddressPort matches the current, there is another disk on the same note, just increase the count
		if _, ok := diskRulesZone.Nodes[device.IP]; ok {
			diskRulesZone.Nodes[device.IP].DiskCount++
			continue
		}

		if diskRulesZone.Nodes == nil {
			diskRulesZone.Nodes = make(map[string]*rules.NodeRules)
		}
		weight := device.Weight
		diskRulesZone.Nodes[device.IP] = &rules.NodeRules{
			DiskCount: 1,
			Weight:    &weight,
		}
	}

	return diskRules
}
