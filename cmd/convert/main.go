// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package convertcmd

import (
	"path/filepath"

	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/go-bits/must"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/sapcc/swift-ring-artisan/pkg/builderfile"
	"github.com/sapcc/swift-ring-artisan/pkg/convert"
	"github.com/sapcc/swift-ring-artisan/pkg/misc"
	"github.com/sapcc/swift-ring-artisan/pkg/rules"
)

var (
	baseSize        float64
	builderFilename string
	outputFilename  string
)

// AddCommandTo adds a command to cobra.Command
func AddCommandTo(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "convert",
		Example: "  swift-ring-artisan convert --input <file> --output <file>",
		Short:   "Does an initial conversion from a parsed swift-ring-builder file to a rule file.",
		Long: `Does an initial conversion from a parsed swift-ring-builder file to a rule file.
This file should be edited and anchors & aliases should be added.
Currently this is not possible from within go due to a bug in the yaml library that support anchors & aliases.`,
		Run: run,
	}
	cmd.PersistentFlags().Float64VarP(&baseSize, "size", "s", 6, "Base size to use for the size of the disks. Defaults to 6TB. Should have a suffix like GB or TB.")
	cmd.PersistentFlags().StringVarP(&builderFilename, "builder", "b", "", "Builder file to read and apply the changes to.")
	cmd.PersistentFlags().StringVarP(&outputFilename, "output", "o", "", "Output file to write the rules to.")
	parent.AddCommand(cmd)
}

func run(cmd *cobra.Command, args []string) {
	if builderFilename == "" {
		logg.Fatal("--builder needs to be supplied and cannot be empty")
	}

	ring := builderfile.File(builderFilename)

	diskRules := convert.Convert(ring, baseSize)

	filename := filepath.Base(builderFilename)
	file := map[string]rules.RingRules{filename: diskRules}

	dataYAML := must.Return(yaml.Marshal(file))
	misc.WriteToStdoutOrFile(dataYAML, outputFilename)
}
