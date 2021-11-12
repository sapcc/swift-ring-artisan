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

package main

import (
	"os"
	"strconv"

	"github.com/sapcc/go-bits/logg"
	applycmd "github.com/sapcc/swift-ring-artisan/cmd/apply"
	convertcmd "github.com/sapcc/swift-ring-artisan/cmd/convert"
	"github.com/spf13/cobra"
)

// ParseBool is like strconv.ParseBool() but doesn't return any error.
func ParseBool(str string) bool {
	v, _ := strconv.ParseBool(str)
	return v
}

func main() {
	logg.ShowDebug = ParseBool(os.Getenv("ARTISAN_DEBUG"))

	rootCmd := &cobra.Command{
		Use:   "swift-ring-artisan",
		Short: "Declarative frontend for swift-ring-builder",
		Long:  "swift-ring-artisan is a declarative frontend for swift-ring-builder. This binary also contains a tool to parse the output of swift-ring-builder to a machine readable format.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	applycmd.AddCommandTo(rootCmd)
	convertcmd.AddCommandTo(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		logg.Fatal(err.Error())
	}
}
