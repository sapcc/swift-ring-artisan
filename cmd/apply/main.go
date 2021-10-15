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

package applycmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/swift-ring-artisan/pkg/parse"
	"github.com/sapcc/swift-ring-artisan/pkg/rules"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	outputFormat   string
	inputFilename  string
	outputFilename string
	ruleFilename   string
)

func AddCommandTo(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "apply -i <file> -r <file>",
		Example: "  swift-ring-builder -i swift-ring-builder-output.yaml -r swift-ring-artisan-rules.yaml",
		Short:   "Applies rules to the parsed swift-ring-builder file.",
		Long: `Generates swift-ring-builder commands based on predefined rules which get applied to the parsed output of the swift-ring-builder utility.
		Rebalance needs to be done manually afterwards.`,
		Run: run,
	}
	cmd.PersistentFlags().StringVarP(&inputFilename, "input", "i", "", "Input file from where the parsed data should be read.")
	cmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "", "Output format. Can be either json or yaml.")
	cmd.PersistentFlags().StringVarP(&outputFilename, "output", "o", "", "Output file to write the parsed data to.")
	cmd.PersistentFlags().StringVarP(&ruleFilename, "rule", "r", "", "Rule file to apply to the input data.")
	parent.AddCommand(cmd)
}

func run(cmd *cobra.Command, args []string) {
	if outputFormat != "" && outputFormat != "json" && outputFormat != "yaml" {
		logg.Fatal("format needs to be set to json OR yaml.")
	}

	// read and parse input file
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

	// read and parse rule file
	if ruleFilename == "" {
		logg.Fatal("--rule needs to be supplied and cannot be empty.")
	}
	ruleContent, err := ioutil.ReadFile(ruleFilename)
	if err != nil {
		logg.Fatal(err.Error())
	}
	var ruleData rules.DiskRules
	err = yaml.Unmarshal(ruleContent, &ruleData)
	if err != nil {
		logg.Fatal(fmt.Sprintf("Parsing file failed: %s", err.Error()))
	}

	logg.Debug(fmt.Sprintf("inputData: %+v", inputData))
	logg.Debug(fmt.Sprintf("ruleData: %+v", ruleData))
	parsedData := rules.ApplyRules(inputData, ruleData)
	if err != nil {
		logg.Fatal(err.Error())
	}

	var parsedDataOutput []byte
	// default to yaml
	if outputFormat == "" || outputFormat == "yaml" {
		parsedDataOutput, err = yaml.Marshal(parsedData)
	} else if outputFormat == "json" {
		parsedDataOutput, err = json.Marshal(parsedData)
	}
	if err != nil {
		logg.Fatal(err.Error())
	}

	if outputFilename == "" {
		fmt.Printf("%+v", string(parsedDataOutput))
	} else {
		err := os.WriteFile(outputFilename, parsedDataOutput, 0644)
		if err != nil {
			logg.Fatal(fmt.Sprintf("writing data to %s failed: %s", outputFilename, err.Error()))
		}
	}

	if err != nil {
		os.Exit(1)
	}
}
