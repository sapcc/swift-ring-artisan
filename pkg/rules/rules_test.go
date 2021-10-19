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

package rules

import (
	"testing"

	"github.com/sapcc/go-bits/assert"
	"github.com/sapcc/swift-ring-artisan/pkg/misc"
	"github.com/sapcc/swift-ring-artisan/pkg/parse"
)

func TestApplyRules1(t *testing.T) {
	var input parse.MetaData
	misc.ReadYAML("../../testing/builder-output.yaml", &input)

	var rules DiskRules
	misc.ReadYAML("../../testing/artisan-rules-changes-1.yaml", &rules)

	commandQueue := ApplyRules(input, rules, "/dev/null")

	assert.DeepEqual(t, "parsing", commandQueue, []string{
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-01 --weight 100 166",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-02 --weight 100 166",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-03 --weight 100 166",
	})
}

func TestApplyRules2(t *testing.T) {
	var input parse.MetaData
	misc.ReadYAML("../../testing/builder-output.yaml", &input)

	var rules DiskRules
	misc.ReadYAML("../../testing/artisan-rules-changes-2.yaml", &rules)

	commandQueue := ApplyRules(input, rules, "/dev/null")

	assert.DeepEqual(t, "parsing", commandQueue, []string{
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-01 --weight 100 166",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-02 --weight 100 166",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-03 --weight 100 166",
	})
}
