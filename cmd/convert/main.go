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

package convertcmd

import (
	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/swift-ring-artisan/pkg/builderfile"
	"github.com/sapcc/swift-ring-artisan/pkg/convert"
	"github.com/sapcc/swift-ring-artisan/pkg/misc"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	basePort       uint64
	baseSize       float64
	inputFilename  string
	outputFilename string
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
	cmd.PersistentFlags().Uint64VarP(&basePort, "port", "p", 6001, "Base port to use for all nodes. Defaults to 6001.")
	cmd.PersistentFlags().Float64VarP(&baseSize, "size", "s", 6, "Base size to use for the size of the disks. Defaults to 6TB. Should have a suffix like GB or TB.")
	cmd.PersistentFlags().StringVarP(&inputFilename, "input", "i", "", "Input file from where the parsed data should be read.")
	cmd.PersistentFlags().StringVarP(&outputFilename, "output", "o", "", "Output file to write the rules to.")
	parent.AddCommand(cmd)
}

func run(cmd *cobra.Command, args []string) {
	if inputFilename == "" {
		logg.Fatal("--input needs to be supplied and cannot be empty.")
	}

	var ring builderfile.RingInfo
	misc.ReadYAML(inputFilename, &ring)

	diskRules := convert.Convert(ring, basePort, baseSize)

	dataYAML, err := yaml.Marshal(diskRules)
	if err != nil {
		logg.Fatal(err.Error())
	}
	misc.WriteToStdoutOrFile(dataYAML, outputFilename)
}
