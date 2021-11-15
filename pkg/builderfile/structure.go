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

package builderfile

import (
	"fmt"
	"time"
)

type DeviceInfo struct {
	ID              uint64
	Region          uint64
	Zone            uint64
	IP              string // TODO: remove
	Port            uint64
	ReplicationIP   string `yaml:"replication_ip" mapstructure:"replication_ip"`
	ReplicationPort uint64 `yaml:"replication_port" mapstructure:"replication_port"`
	Name            string `mapstructure:"device"` // TODO: rename to Device?
	Weight          float64
	Partitions      uint64 `mapstructure:"parts"`
	Balance         float64
	//lint:ignore U1000 TODO
	flags struct{} // TODO: figure out how the field looks like
	//lint:ignore U1000 TODO
	meta struct{} // TODO: figure out how the field looks like
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
	return fmt.Sprintf("%s:%d", device.IP, device.Port)
}

func (device DeviceInfo) CommandAdd(ringFilename string) string {
	return fmt.Sprintf("swift-ring-builder %s add --region %d --zone %d --ip %s --port %d --device %s --weight %g",
		ringFilename, device.Region, device.Zone, device.IP, device.Port, device.Name, device.Weight)
}

func (device DeviceInfo) CommandSetWeight(ringFilename string, desiredWeight float64) string {
	return fmt.Sprintf("swift-ring-builder %s set_weight --region %d --zone %d --ip %s --port %d --device %s --weight %g %g",
		ringFilename, device.Region, device.Zone, device.IP, device.Port, device.Name, device.Weight, desiredWeight)
}

func (device DeviceInfo) CommandRemove(ringFilename string) string {
	return fmt.Sprintf("swift-ring-builder %s remove --region %d --zone %d --ip %s --port %d --device %s --weight %g",
		ringFilename, device.Region, device.Zone, device.IP, device.Port, device.Name, device.Weight)
}
