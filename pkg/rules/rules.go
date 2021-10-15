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
	"github.com/sapcc/swift-ring-artisan/pkg/parse"
)

// Disk properties
type Disk struct {
	Count  uint64
	Size   string
	Weight float64
}

// Node is a server containing disks
type Node struct {
	IPPort string   `yaml:"ip_port"`
	Meta   struct{} `yaml:"meta,omitempty"` // TODO: figure out how the field looks like
	Disks  Disk
}

// Zone contains multiple nodes
type Zone struct {
	ID    uint64
	Nodes []Node
}

// DiskRules containing the rules for a region, multiple Zones and dozzens Nodes
type DiskRules struct {
	BaseSize string `yaml:"base_size"`
	Region   uint64
	Zones    []Zone
}

// ApplyRules to parsed MetaData
func ApplyRules(inputData parse.MetaData, ruleData DiskRules) parse.MetaData {
	return inputData
}
