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

package parsecmd

import (
	"fmt"
	"io"
	"os"

	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/go-bits/must"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/sapcc/swift-ring-artisan/pkg/builderfile"
)

var (
	outputFile string
)

// AddCommandTo adds a command to cobra.Command
func AddCommandTo(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "parse <file>",
		Example: "  swift-ring-builder object.builder | swift-ring-artisan parse",
		Short:   "Parses swift-ring-builder output into yaml. Mainly used for test cases.",
		Long: `Parses swift-ring-builder output into yaml. Mainly used for test cases.
The output of swift-ring-builder needs to be piped into swift-ring-artisan.`,
		Run: run,
	}
	cmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "Output file to write the parsed data to.")
	parent.AddCommand(cmd)
}

func run(cmd *cobra.Command, args []string) {
	// if a file is supplied read that otherwise listen on stdin
	var input io.Reader
	if len(args) == 1 {
		file, err := os.Open(args[0])
		if err != nil {
			logg.Fatal("Reading file failed: %s", err.Error())
		}
		defer file.Close()
		input = file
	} else {
		input = os.Stdin
	}

	metaData := builderfile.Input(input)
	metaDataOutput := must.Return(yaml.Marshal(metaData))

	if outputFile == "" {
		fmt.Printf("%+v", string(metaDataOutput))
	} else {
		err := os.WriteFile(outputFile, metaDataOutput, 0644)
		if err != nil {
			logg.Fatal("writing data to %s failed: %s", outputFile, err.Error())
		}
	}
}
