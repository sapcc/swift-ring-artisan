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
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/swift-ring-artisan/pkg/convert"
	"github.com/sapcc/swift-ring-artisan/pkg/misc"
	"github.com/sapcc/swift-ring-artisan/pkg/parse"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	baseSize       string
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
	cmd.PersistentFlags().StringVarP(&baseSize, "size", "s", "6TB", "Base size to use for the size of the disks. Defaults to 6TB. Should have a suffix like GB or TB.")
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

	diskRules := convert.Convert(inputData, baseSize)

	dataYAML, err := yaml.Marshal(diskRules)
	if err != nil {
		logg.Fatal(err.Error())
	}

	misc.WriteToStdoutOrFile(dataYAML, outputFilename)

	if err != nil {
		os.Exit(1)
	}
}
