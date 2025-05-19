// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"strconv"

	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/go-bits/must"
	"github.com/spf13/cobra"

	applycmd "github.com/sapcc/swift-ring-artisan/cmd/apply"
	convertcmd "github.com/sapcc/swift-ring-artisan/cmd/convert"
	parsecmd "github.com/sapcc/swift-ring-artisan/cmd/parse"
)

// ParseBool is like strconv.ParseBool() but doesn't return any error.
func ParseBool(str string) bool {
	v, _ := strconv.ParseBool(str) //nolint:errcheck
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
			cmd.Help() //nolint:errcheck
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&logg.ShowDebug, "debug", "d", false, "Enable debug log")

	applycmd.AddCommandTo(rootCmd)
	convertcmd.AddCommandTo(rootCmd)
	parsecmd.AddCommandTo(rootCmd)

	must.Succeed(rootCmd.Execute())
}
