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

package initcmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/swift-ring-artisan/pkg/parse"
	"github.com/sapcc/swift-ring-artisan/pkg/rules"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	inputFilename  string
	outputFilename string
)

func AddCommandTo(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "init",
		Example: "  swift-ring-artisan init --input <file> --output <file>",
		Short:   "Does an initial conversion from a parsed swift-ring-builder file to a rule file.",
		Long: `Does an initial conversion from a parsed swift-ring-builder file to a rule file.
This file should be edited and anchors & aliases should be added.
Currently this is not possible from within go due to a bug in the yaml library that support anchors & aliases.`,
		Run: run,
	}
	cmd.PersistentFlags().StringVarP(&inputFilename, "input", "i", "", "Input file from where the parsed data should be read.")
	cmd.PersistentFlags().StringVarP(&outputFilename, "output", "o", "", "Output file to write the rules to.")
	parent.AddCommand(cmd)
}

func run(cmd *cobra.Command, args []string) {
	if inputFilename == "" {
		logg.Fatal("--input needs to be supplied and cannot be empty.")
	}

	inputContent, err := ioutil.ReadFile(inputFilename)
	if err != nil {
		logg.Fatal(err.Error())
	}

	var inputData parse.MetaData
	err = yaml.Unmarshal(inputContent, &inputData)
	if err != nil {
		logg.Fatal(fmt.Sprintf("Parsing file failed: %s", err.Error()))
	}

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
		Region: inputData.Devices[0].Region,
		// TODO: better default? argument?
		BaseSize: "6TB",
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

	dataYAML, err := yaml.Marshal(diskRules)
	if err != nil {
		logg.Fatal(err.Error())
	}

	if outputFilename == "" {
		fmt.Printf("%+v", string(dataYAML))
	} else {
		err := os.WriteFile(outputFilename, dataYAML, 0644)
		if err != nil {
			logg.Fatal(fmt.Sprintf("writing data to %s failed: %s", dataYAML, err.Error()))
		}
	}

	if err != nil {
		os.Exit(1)
	}
}
