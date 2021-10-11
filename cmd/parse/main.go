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
	artisan "github.com/sapcc/swift-ring-artisan/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	format string
)

func AddCommandTo(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "parse <image>",
		Example: "  swift-ring-builder object.builder | swift-ring-artisan parse",
		Short:   "Parses the output of swift-ring-builder to a machine readable format.",
		Long: `Parses the output of swift-ring-builder to a machine readable format.
The output of swift-ring-builder needs to be piped into swift-ring-artisan.`,
		Run: run,
	}
	cmd.PersistentFlags().StringVar(&format, "format", "[]", "Output format. Can be either json or yaml.")
	parent.AddCommand(cmd)
}

func run(cmd *cobra.Command, args []string) {
	// if a file is supplied read that otherwise listen on stdin
	var input io.Reader
	if len(args) == 1 {
		file, err := os.Open(args[0])
		if err != nil {
			logg.Fatal(fmt.Sprintf("Reading file failed: %s", err.Error()))
		}
		defer file.Close()
		input = file
	} else {
		input = os.Stdin
	}

	metaData, err := artisan.ParseInput(input)
	if err != nil {
		logg.Fatal(err.Error())
	}

	metaDataYaml, err := yaml.Marshal(metaData)
	if err != nil {
		logg.Fatal(err.Error())
	}

	fmt.Printf("%+v", string(metaDataYaml))

	if err != nil {
		os.Exit(1)
	}
}
