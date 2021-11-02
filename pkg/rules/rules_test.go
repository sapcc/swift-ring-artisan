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
	"github.com/sapcc/swift-ring-artisan/pkg/builderfile"
	"github.com/sapcc/swift-ring-artisan/pkg/misc"
)

func TestApplyRules1(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-1.yaml", &input)

	var rules RingRules
	misc.ReadYAML("../../testing/artisan-rules-changes-1.yaml", &rules)

	commandQueue := rules.CalculateChanges(input, "/dev/null")

	assert.DeepEqual(t, "parsing", commandQueue, []string{
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-01 --weight 100 166",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-02 --weight 100 166",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-03 --weight 100 166",
	})
}

func TestApplyRules1_1(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-1.yaml", &input)

	var rules RingRules
	misc.ReadYAML("../../testing/artisan-rules-changes-1-1.yaml", &rules)

	commandQueue := rules.CalculateChanges(input, "/dev/null")

	assert.DeepEqual(t, "parsing", commandQueue, []string{
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-01 --weight 100 166",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-02 --weight 100 166",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-03 --weight 100 166",
	})
}

func TestApplyRules2(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-2.yaml", &input)

	var rules RingRules
	misc.ReadYAML("../../testing/artisan-rules-changes-2.yaml", &rules)

	commandQueue := rules.CalculateChanges(input, "/dev/null")

	assert.DeepEqual(t, "parsing", commandQueue, []string{
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-01 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-02 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-03 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-04 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-05 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-06 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-07 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-08 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-09 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-10 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-11 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-12 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-13 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-14 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-15 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-16 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-17 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-18 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-19 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-20 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-21 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-22 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-23 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-24 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-25 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-26 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-27 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-28 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-29 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-30 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-31 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-32 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-33 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-34 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-35 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-36 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-37 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-38 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-39 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-40 --weight 100 200",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-01 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-02 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-03 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-04 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-05 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-06 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-07 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-08 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-09 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-10 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-11 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-12 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-13 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-14 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-15 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-16 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-17 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-18 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-19 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-20 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-21 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-22 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-23 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-24 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-25 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-26 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-27 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-28 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-29 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-30 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-31 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-32 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-33 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-34 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-35 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-36 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-37 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-38 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-39 --weight 100 150",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-40 --weight 100 150",
	})
}
