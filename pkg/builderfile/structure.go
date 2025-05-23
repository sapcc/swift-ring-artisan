// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package builderfile

import (
	"encoding/json"
	"fmt"
	"time"
)

type DeviceInfo struct {
	ID              uint64
	Region          uint64
	Zone            uint64
	NodeIP          string `yaml:"ip" mapstructure:"ip"`
	Port            uint64
	ReplicationIP   string `yaml:"replication_ip" mapstructure:"replication_ip"`
	ReplicationPort uint64 `yaml:"replication_port" mapstructure:"replication_port"`
	Name            string `mapstructure:"device"`
	Weight          float64
	Partitions      uint64 `mapstructure:"parts"`
	Balance         float64
	Meta            *map[string]string `yaml:"meta,omitempty"`
	//nolint:unused
	flags struct{} // TODO: figure out how the field looks like
}

// RingInfo contains the meta data about the ring file
type RingInfo struct {
	FileName string `yaml:"file_name"`
	Version  uint64
	ID       string

	Partitions  uint64
	Replicas    float64
	Regions     uint64
	Zones       uint64
	DeviceCount uint64 `yaml:"device_count"`
	Balance     float64
	Dispersion  float64

	ReassignedCooldown  uint64    `yaml:"reassigned_cooldown"`
	ReassignedRemaining time.Time `yaml:"reassigned_remaining"`

	OverloadFactorPercent float64 `yaml:"overload_factor_Percent"`
	OverloadFactorDecimal float64 `yaml:"overload_factor_decimal"`

	Devices []DeviceInfo
}

func (device DeviceInfo) IPAddressPort() string {
	return fmt.Sprintf("%s:%d", device.NodeIP, device.Port)
}

func (ring RingInfo) CommandSetOverload(ringFilename string, desiredOverload float64) string {
	return fmt.Sprintf("swift-ring-builder %s set_overload %f", ringFilename, desiredOverload)
}

func (device DeviceInfo) CommandAdd(ringFilename string) string {
	if device.Meta != nil {
		var meta []byte
		meta, _ = json.Marshal(device.Meta) //nolint:errcheck
		return fmt.Sprintf("swift-ring-builder %s add --region %d --zone %d --ip %s --port %d --device %s --weight %g --meta %s",
			ringFilename, device.Region, device.Zone, device.NodeIP, device.Port, device.Name, device.Weight, meta)
	}

	return fmt.Sprintf("swift-ring-builder %s add --region %d --zone %d --ip %s --port %d --device %s --weight %g",
		ringFilename, device.Region, device.Zone, device.NodeIP, device.Port, device.Name, device.Weight)
}

func (device DeviceInfo) CommandSetMeta(ringFilename string, desiredMeta map[string]string) string {
	var meta []byte
	meta, _ = json.Marshal(desiredMeta) //nolint:errcheck
	return fmt.Sprintf("swift-ring-builder %s set_info --region %d --zone %d --ip %s --port %d --device %s --change-meta %s",
		ringFilename, device.Region, device.Zone, device.NodeIP, device.Port, device.Name, meta)
}

func (device DeviceInfo) CommandSetMetaNode(ringFilename string, desiredMeta map[string]string) string {
	var meta []byte
	meta, _ = json.Marshal(desiredMeta) //nolint:errcheck
	return fmt.Sprintf("swift-ring-builder %s set_info --region %d --zone %d --ip %s --port %d --change-meta %s --yes",
		ringFilename, device.Region, device.Zone, device.NodeIP, device.Port, meta)
}

func (device DeviceInfo) CommandSetWeight(ringFilename string, desiredWeight float64) string {
	return fmt.Sprintf("swift-ring-builder %s set_weight --region %d --zone %d --ip %s --port %d --device %s --weight %g %g",
		ringFilename, device.Region, device.Zone, device.NodeIP, device.Port, device.Name, device.Weight, desiredWeight)
}

func (device DeviceInfo) CommandRemove(ringFilename string) string {
	return fmt.Sprintf("swift-ring-builder %s remove --region %d --zone %d --ip %s --port %d --device %s --weight %g",
		ringFilename, device.Region, device.Zone, device.NodeIP, device.Port, device.Name, device.Weight)
}
