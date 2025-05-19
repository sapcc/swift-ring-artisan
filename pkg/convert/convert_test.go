// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package convert

import (
	"testing"

	"github.com/sapcc/go-bits/assert"

	"github.com/sapcc/swift-ring-artisan/pkg/builderfile"
	"github.com/sapcc/swift-ring-artisan/pkg/misc"
	"github.com/sapcc/swift-ring-artisan/pkg/rules"
)

func TestParse1(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-1.yaml", &input)

	var expected rules.RingRules
	misc.ReadYAML("../../testing/artisan-rules-1.yaml", &expected)

	metaData := Convert(input, 6)
	assert.DeepEqual(t, "parsing", metaData, expected)
}

func TestParse2(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-2.yaml", &input)

	var expected rules.RingRules
	misc.ReadYAML("../../testing/artisan-rules-2.yaml", &expected)

	metaData := Convert(input, 6)
	assert.DeepEqual(t, "parsing", metaData, expected)
}
