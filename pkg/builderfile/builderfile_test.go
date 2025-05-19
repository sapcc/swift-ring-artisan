// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package builderfile

import (
	"os"
	"testing"

	"github.com/sapcc/go-bits/assert"
	"github.com/sapcc/go-bits/logg"

	"github.com/sapcc/swift-ring-artisan/pkg/misc"
)

func TestParse1(t *testing.T) {
	input, err := os.Open("../../testing/builder-output-1.txt")
	if err != nil {
		logg.Fatal("Reading file failed: %s", err.Error())
	}
	defer input.Close()

	var expected RingInfo
	misc.ReadYAML("../../testing/builder-output-1.yaml", &expected)

	metaData := Input(input)
	assert.DeepEqual(t, "parsing", metaData, expected)
}

func TestParse2(t *testing.T) {
	input, err := os.Open("../../testing/builder-output-2.txt")
	if err != nil {
		logg.Fatal("Reading file failed: %s", err.Error())
	}
	defer input.Close()

	var expected RingInfo
	misc.ReadYAML("../../testing/builder-output-2.yaml", &expected)

	metaData := Input(input)
	assert.DeepEqual(t, "parsing", metaData, expected)
}
