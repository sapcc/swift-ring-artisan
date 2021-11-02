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
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/swift-ring-artisan/pkg/builderfile"
	"github.com/sapcc/swift-ring-artisan/pkg/misc"
	"github.com/sapcc/swift-ring-artisan/pkg/pickler"
	"github.com/sapcc/swift-ring-artisan/pkg/rules"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

var (
	checkChanges    bool
	executeCommands bool
	inputFilename   string
	outputFilename  string
	outputFormat    string
	builderFilename string
	ruleFilename    string
)

// AddCommandTo adds a command to cobra.Command
func AddCommandTo(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "apply -i <file> -r <file>",
		Example: "  swift-ring-artisan apply -b account.builder -r swift-ring-artisan-rules.yaml",
		Short:   "Applies rules to the parsed swift-ring-builder file.",
		Long: `Generates swift-ring-builder commands based on predefined rules which get applied to the parsed output of the swift-ring-builder utility.
		Rebalance needs to be done manually afterwards.`,
		Run: run,
	}
	cmd.PersistentFlags().BoolVarP(&checkChanges, "check", "c", false, "Wether to check if the rule file matches the ring. If it does not match the exit code is 1.")
	cmd.PersistentFlags().BoolVarP(&executeCommands, "execute", "e", false, "Wether to execute the generated commands.")
	cmd.PersistentFlags().StringVarP(&inputFilename, "input", "i", "", "Input file from where the parsed data should be read.")
	cmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "", "Output format. Can be either json or yaml.")
	cmd.PersistentFlags().StringVarP(&outputFilename, "output", "o", "", "Output file to write the parsed data to.")
	cmd.PersistentFlags().StringVarP(&builderFilename, "builder", "b", "", "Builder file to apply the changes to.")
	cmd.PersistentFlags().StringVarP(&ruleFilename, "rule", "r", "", "Rule file to apply to the input data.")
	parent.AddCommand(cmd)
}

func run(cmd *cobra.Command, args []string) {
	if outputFormat != "" && outputFormat != "json" && outputFormat != "yaml" {
		logg.Fatal("format needs to be set to json OR yaml.")
	}

	var ring builderfile.RingInfo
	if inputFilename == "" && builderFilename != "" {
		// // generate with ./unpickle.sh
		// cmd := exec.Command("python3", "-c", "'import json;import pickle;import sys;d=pickle.load(open(sys.argv[-1],\"rb\"));d[\"_dispersion_graph\"]=None;d[\"_replica2part2dev\"]=None;d[\"_last_part_moves\"]=None;print(json.dumps(d));'", builderFilename)
		// stdout, err := cmd.Output()
		// logg.Info("%s %s", cmd, string(stdout))
		// if err == exec.ErrNotFound {
		// 	logg.Fatal("Please install python3 to decode the pickle file.")
		// } else if err != nil {
		// 	logg.Fatal(err.Error())
		// }

		// var pickleData PickleData
		// decoder := json.NewDecoder(bytes.NewReader(stdout))
		// decoder.DisallowUnknownFields()
		// err = decoder.Decode(&pickleData)
		// if err != nil {
		// 	logg.Fatal(err.Error())
		// }

		pickleData := pickler.DecodeBuilderFile(builderFilename)

		// optional compare pickle parser with cli parser
		cmd := exec.Command("sh", "-c", "command -v swift-ring-builder")
		err := cmd.Run()
		if err == exec.ErrNotFound {
			logg.Debug("Did not find swift-ring-builder, skipping consistency check")
		} else if err == nil {
			cmd := exec.Command("swift-ring-builder", builderFilename)
			stdout, err := cmd.Output()
			if err != nil {
				logg.Fatal(err.Error())
			}
			ring = builderfile.Input(bytes.NewReader(stdout))

			equal := reflect.DeepEqual(ring, pickleData)
			if !equal {
				dmp := diffmatchpatch.New()
				diffs := dmp.DiffMain(fmt.Sprintf("%+v\n", ring), fmt.Sprintf("%+v\n", pickleData), false)
				logg.Info("Pickle parsed output and swift-ring-builder output are not equal. What is going on here?!")
				logg.Fatal(dmp.DiffPrettyText(diffs))
			}
		} else {
			logg.Fatal(err.Error())
		}

		ring = builderfile.RingInfo{
			ID:         pickleData.ID,
			Replicas:   pickleData.Replicas,
			Regions:    pickleData.Devs[0].Region, // FIXME: make multi region aware
			Dispersion: pickleData.Dispersion,
			Devices:    pickleData.Devs,
		}
	} else if inputFilename != "" {
		misc.ReadYAML(inputFilename, &ring)
	} else {
		logg.Fatal("Either --input or --builder needs to be set.")
	}

	if ruleFilename == "" {
		logg.Fatal("--rule needs to be supplied and cannot be empty.")
	}
	var ruleData rules.RingRules
	misc.ReadYAML(ruleFilename, &ruleData)

	if (checkChanges || executeCommands) && builderFilename == "" {
		logg.Fatal("--ring needs to be supplied and cannot be empty.")
	}
	commandQueue := ruleData.CalculateChanges(ring, builderFilename)
	if checkChanges {
		if len(commandQueue) != 0 {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}

	for _, command := range commandQueue {
		if executeCommands {
			args := strings.Split(command, " ")
			cmd := exec.Command(args[0], args[1:]...)
			stdout, err := cmd.Output()
			if err != nil {
				logg.Fatal("Command \"%s\" failed: %v", command, err.Error())
			}
			logg.Info(string(stdout))
		} else {
			misc.WriteToStdoutOrFile([]byte(command+"\n"), outputFilename)
		}
	}
}
