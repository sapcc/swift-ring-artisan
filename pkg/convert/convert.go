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
	"sort"

	"github.com/sapcc/swift-ring-artisan/pkg/parse"
	"github.com/sapcc/swift-ring-artisan/pkg/rules"
)

func Convert(inputData parse.MetaData, baseSize string) rules.DiskRules {
	sort.Slice(inputData.Devices, func(i, j int) bool {
		lhs := inputData.Devices[i]
		rhs := inputData.Devices[j]
		if lhs.Zone != rhs.Zone {
			return lhs.Zone < rhs.Zone
		}
		if lhs.IPAddressPort != rhs.IPAddressPort {
			return lhs.IPAddressPort < rhs.IPAddressPort
		}
		return lhs.Name < rhs.Name
	})

	diskRules := rules.DiskRules{
		Region:   inputData.Devices[0].Region,
		BaseSize: baseSize,
	}

	var (
		last          string
		diskRulesZone *rules.Zone
	)
	for _, device := range inputData.Devices {
		// create zone if it does not exist
		if len(diskRules.Zones) == 0 || diskRules.Zones[len(diskRules.Zones)-1].ID != device.Zone {
			diskRules.Zones = append(diskRules.Zones, rules.Zone{ID: device.Zone})
		}

		diskRulesZone = &diskRules.Zones[device.Zone-1]
		// if the last IPAddressPort matches the current, there is another disk on the same note, just increase the count
		if last != "" && last == device.IPAddressPort {
			diskRulesZone.Nodes[len(diskRulesZone.Nodes)-1].Disks.Count += 1
			continue
		}
		diskRulesZone.Nodes = append(diskRulesZone.Nodes, rules.Node{
			IPPort: device.IPAddressPort,
			Disks: rules.Disk{
				Count:  1,
				Size:   diskRules.BaseSize,
				Weight: device.Weight,
			},
		})

		last = device.IPAddressPort
	}

	return diskRules
}
