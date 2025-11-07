// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package applycmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sapcc/go-bits/errext"
	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/go-bits/must"
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

func run(cmd *cobra.Command, args []string) {
	_, _ = cmd, args

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
	ringRules, ok := file[builderBaseFilename]
	if !ok {
		logg.Fatal("%s is missing key for %s", ruleFilename, builderBaseFilename)
	}

	commandQueue, confirmations, err := ringRules.CalculateChanges(ring, builderFilename)
	if err != nil {
		logg.Fatal(err.Error())
	}
	if len(commandQueue) == 0 {
		os.Exit(0)
	}

	if executeCommands && checkChanges {
		logg.Fatal("Cannot execute commands and check if builder and ring file matches.")
	}

	if len(confirmations) > 0 {
		for _, confirmation := range confirmations {
			logg.Info(confirmation)
		}

		if !misc.Prompt("Please type upper-case YES to continue.", []string{"YES"}) {
			logg.Fatal("Aborting")
		}
	}

	if !executeCommands {
		for _, command := range commandQueue {
			misc.WriteToStdoutOrFile([]byte(command+"\n"), outputFilename)
		}
	}

	// exit early when only checking for changes to skip executing commands
	if checkChanges {
		if len(commandQueue) > 0 {
			os.Exit(1)
		}
		os.Exit(0)
	}

	promptAnswer := false
	fileInfo := must.Return(os.Stdin.Stat())

	// evaluates to true if program is run in an interactive shell and not piped
	isInteractive := (fileInfo.Mode() & os.ModeCharDevice) != 0
	if !executeCommands && isInteractive {
		promptAnswer = misc.AskConfirmation("Do you want to apply the changes by executing the above commands?")
	}

	rebalanceRequired := false
	if executeCommands || promptAnswer {
		for _, command := range commandQueue {
			// rebalance not required, if commandQueue only contains 'set_info' commands
			rebalanceRequired = rebalanceRequired || !strings.Contains(command, "set_info")
			args := strings.Split(command, " ")
			cmd := exec.Command(args[0], args[1:]...) //nolint:gosec // input is user supplied and self executed
			stdout, err := cmd.Output()
			logg.Info(string(stdout))
			if err != nil {
				logg.Fatal("Command %q failed: %v", command, err.Error())
			}
		}
	} else {
		os.Exit(1)
	}

	promptAnswer = false
	action := "write_ring"
	if rebalanceRequired {
		action = "rebalance"
	}
	if !executeCommands && isInteractive {
		promptAnswer = misc.AskConfirmation(fmt.Sprintf("Do you want to %s now?", action))
	}

	if executeCommands || promptAnswer {
		cmd := exec.Command("swift-ring-builder", builderFilename, action)
		logg.Info(fmt.Sprintf("%s %s", builderFilename, action))
		stdout, err := cmd.Output()
		// For better readablitity, split multiline outputs to separate loglines
		for line := range strings.SplitSeq(string(stdout), "\n") {
			if line != "" {
				logg.Info(line)
			}
		}

		if exitError, ok := errext.As[*exec.ExitError](err); ok {
			os.Exit(exitError.ExitCode())
		} else if err != nil {
			logg.Fatal("Command %q failed: %v", strings.Join(cmd.Args, " "), err.Error())
		}
	}

	os.Exit(0)
}
