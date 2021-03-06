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
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sapcc/go-bits/logg"
	"github.com/spf13/cobra"

	"github.com/sapcc/swift-ring-artisan/pkg/builderfile"
	"github.com/sapcc/swift-ring-artisan/pkg/misc"
	"github.com/sapcc/swift-ring-artisan/pkg/rules"
)

var (
	checkChanges    bool
	executeCommands bool
	outputFilename  string
	outputFormat    string
	builderFilename string
	ruleFilename    string
)

// AddCommandTo adds a command to cobra.Command
func AddCommandTo(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "apply -b <file> -r <file>",
		Example: "  swift-ring-artisan apply -b account.builder -r swift-ring-artisan-rules.yaml",
		Short:   "Applies rules to a swift-ring-builder file.",
		Long: `Generates swift-ring-builder commands based on predefined rules which get applied to the parsed output of the swift-ring-builder utility.
		Rebalance needs to be done manually afterwards.`,
		Run: run,
	}
	cmd.PersistentFlags().BoolVarP(&checkChanges, "check", "c", false, "Wether to check if the rule file matches the ring. If it does not match the exit code is 1.")
	cmd.PersistentFlags().BoolVarP(&executeCommands, "execute", "e", false, "Wether to execute the generated commands.")
	cmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "", "Output format. Can be either json or yaml.")
	cmd.PersistentFlags().StringVarP(&outputFilename, "output", "o", "", "Output file to write the parsed data to.")
	cmd.PersistentFlags().StringVarP(&builderFilename, "builder", "b", "", "Builder file to read and apply the changes to.")
	cmd.PersistentFlags().StringVarP(&ruleFilename, "rule", "r", "", "Rule file to apply to the input data.")
	parent.AddCommand(cmd)
}

func run(_ *cobra.Command, args []string) {
	if outputFormat != "" && outputFormat != "json" && outputFormat != "yaml" {
		logg.Fatal("format needs to be set to json OR yaml.")
	}
	if builderFilename == "" {
		logg.Fatal("--builder needs to be set")
	}
	ring := builderfile.File(builderFilename)

	if ruleFilename == "" {
		logg.Fatal("--rule needs to be supplied and cannot be empty")
	}
	var file map[string]rules.RingRules
	misc.ReadYAML(ruleFilename, &file)

	builderBaseFilename := filepath.Base(builderFilename)
	rules, ok := file[builderBaseFilename]
	if !ok {
		logg.Fatal("%s is missing key for %s", ruleFilename, builderBaseFilename)
	}

	commandQueue, err := rules.CalculateChanges(ring, builderFilename)
	if err != nil {
		log.Fatal(err.Error())
	}
	if len(commandQueue) == 0 {
		os.Exit(0)
	}

	if executeCommands && checkChanges {
		logg.Fatal("Cannot execute commands and check if builder and ring file matches.")
	}
	for _, command := range commandQueue {
		misc.WriteToStdoutOrFile([]byte(command+"\n"), outputFilename)
	}

	// exit early when only checking for changes to skip executing commands
	if checkChanges {
		os.Exit(0)
	}

	promptAnswer := false
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err.Error())
	}

	// evaluates to true if program is run in an interactive shell and not piped
	isInteractive := (fileInfo.Mode() & os.ModeCharDevice) != 0
	if isInteractive {
		promptAnswer = misc.AskConfirmation("Do you want to apply the changes by executing the above commands?")
	}

	rebalanceRequired := false
	if executeCommands || promptAnswer {
		for _, command := range commandQueue {
			// rebalance not required, if commandQueue only contains 'set_info' commands
			rebalanceRequired = rebalanceRequired || !strings.Contains(command, "set_info")
			args := strings.Split(command, " ")
			cmd := exec.Command(args[0], args[1:]...)
			stdout, err := cmd.Output()
			logg.Info(string(stdout))
			if err != nil {
				logg.Fatal("Command %q failed: %v", command, err.Error())
			}
		}
	} else {
		os.Exit(1)
	}

	if isInteractive {
		promptAnswer := false
		action := "write_ring"
		if rebalanceRequired {
			action = "rebalance"
		}
		promptAnswer = misc.AskConfirmation(fmt.Sprintf("Do you want to %s now?", action))

		if promptAnswer {
			cmd := exec.Command("swift-ring-builder", builderFilename, action)
			stdout, err := cmd.Output()
			logg.Info(string(stdout))
			if exitError, ok := err.(*exec.ExitError); ok {
				os.Exit(exitError.ExitCode())
			} else if err != nil {
				logg.Fatal("Command %q failed: %v", strings.Join(cmd.Args, " "), err.Error())
			}
		}
	}
	os.Exit(0)
}
