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

package convert

import (
	"testing"

	"github.com/sapcc/go-bits/assert"
	"github.com/sapcc/swift-ring-artisan/pkg/misc"
	"github.com/sapcc/swift-ring-artisan/pkg/parse"
	"github.com/sapcc/swift-ring-artisan/pkg/rules"
)

func TestParse1(t *testing.T) {
	var input parse.MetaData
	misc.ReadYAML("../../testing/builder-output-1.yaml", &input)

	var expected rules.DiskRules
	misc.ReadYAML("../../testing/artisan-rules-1.yaml", &expected)

	metaData := Convert(input, 6001, 6)
	assert.DeepEqual(t, "parsing", metaData, expected)
}

func TestParse2(t *testing.T) {
	var input parse.MetaData
	misc.ReadYAML("../../testing/builder-output-2.yaml", &input)

	var expected rules.DiskRules
	misc.ReadYAML("../../testing/artisan-rules-2.yaml", &expected)

	metaData := Convert(input, 6001, 6)
	assert.DeepEqual(t, "parsing", metaData, expected)
}
