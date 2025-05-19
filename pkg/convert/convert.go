// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

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
		if _, ok := diskRulesZone.Nodes[device.NodeIP]; ok {
			diskRulesZone.Nodes[device.NodeIP].DiskCount++
			continue
		}

		if diskRulesZone.Nodes == nil {
			diskRulesZone.Nodes = make(map[string]*rules.NodeRules)
		}
		weight := device.Weight
		diskRulesZone.Nodes[device.NodeIP] = &rules.NodeRules{
			DiskCount: 1,
			Weight:    &weight,
		}
	}

	return diskRules
}
