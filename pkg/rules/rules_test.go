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
	"fmt"
	"log"
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

	commandQueue, err := rules.CalculateChanges(input, "/dev/null")
	if err != nil {
		log.Fatal(err.Error())
	}

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

	commandQueue, err := rules.CalculateChanges(input, "/dev/null")
	if err != nil {
		log.Fatal(err.Error())
	}

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

	commandQueue, err := rules.CalculateChanges(input, "/dev/null")
	if err != nil {
		log.Fatal(err.Error())
	}

	var expectedCommands []string
	for i := 1; i <= 40; i++ {
		expectedCommands = append(expectedCommands, fmt.Sprintf("swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-%02d --weight 100 200", i))
	}
	for i := 1; i <= 40; i++ {
		expectedCommands = append(expectedCommands, fmt.Sprintf("swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-%02d --weight 100 150", i))
	}
	assert.DeepEqual(t, "parsing", commandQueue, expectedCommands)
}

func TestAddDisk1(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-1.yaml", &input)

	var rules RingRules
	misc.ReadYAML("../../testing/artisan-addition-1.yaml", &rules)

	commandQueue, err := rules.CalculateChanges(input, "/dev/null")
	if err != nil {
		log.Fatal(err.Error())
	}

	assert.DeepEqual(t, "parsing", commandQueue, []string{
		"swift-ring-builder /dev/null add --region 1 --zone 1 --ip 10.114.1.204 --port 6001 --device swift-01 --weight 100",
		"swift-ring-builder /dev/null add --region 1 --zone 1 --ip 10.114.1.204 --port 6001 --device swift-02 --weight 100",
		"swift-ring-builder /dev/null add --region 1 --zone 1 --ip 10.114.1.204 --port 6001 --device swift-03 --weight 100",
	})
}

func TestAddDisk2(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-2.yaml", &input)

	var rules RingRules
	misc.ReadYAML("../../testing/artisan-addition-2.yaml", &rules)

	commandQueue, err := rules.CalculateChanges(input, "/dev/null")
	if err != nil {
		log.Fatal(err.Error())
	}

	var expectedCommands []string
	for i := 1; i <= 12; i++ {
		expectedCommands = append(expectedCommands, fmt.Sprintf("swift-ring-builder /dev/null add --region 1 --zone 4 --ip 10.46.14.161 --port 6001 --device swift-%02d --weight 166", i))
	}
	for i := 1; i <= 12; i++ {
		expectedCommands = append(expectedCommands, fmt.Sprintf("swift-ring-builder /dev/null add --region 1 --zone 4 --ip 10.46.14.248 --port 6001 --device swift-%02d --weight 166", i))
	}
	for i := 1; i <= 12; i++ {
		expectedCommands = append(expectedCommands, fmt.Sprintf("swift-ring-builder /dev/null add --region 1 --zone 4 --ip 10.46.14.42 --port 6001 --device swift-%02d --weight 166", i))
	}
	assert.DeepEqual(t, "parsing", commandQueue, expectedCommands)
}

func TestDeleteDisk1(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-1.yaml", &input)

	var rules RingRules
	misc.ReadYAML("../../testing/artisan-deletion-1.yaml", &rules)

	commandQueue, err := rules.CalculateChanges(input, "/dev/null")
	if err != nil {
		log.Fatal(err.Error())
	}

	assert.DeepEqual(t, "parsing", commandQueue, []string{
		"swift-ring-builder /dev/null remove --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-01 --weight 100",
		"swift-ring-builder /dev/null remove --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-02 --weight 100",
		"swift-ring-builder /dev/null remove --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-03 --weight 100",
	})
}

func TestDeleteDisk2(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-2.yaml", &input)

	var rules RingRules
	misc.ReadYAML("../../testing/artisan-deletion-2.yaml", &rules)

	commandQueue, err := rules.CalculateChanges(input, "/dev/null")
	if err != nil {
		log.Fatal(err.Error())
	}

	var expectedCommands []string
	for i := 1; i <= 40; i++ {
		expectedCommands = append(expectedCommands, fmt.Sprintf("swift-ring-builder /dev/null remove --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-%02d --weight 100", i))
	}
	assert.DeepEqual(t, "parsing", commandQueue, expectedCommands)
}

func TestZoneMismatch(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-1.yaml", &input)

	var rules RingRules
	misc.ReadYAML("../../testing/artisan-error-zones.yaml", &rules)

	_, err := rules.CalculateChanges(input, "/dev/null")
	if err == nil {
		t.Fatal("This test is expected to fail")
	}
	errString := "zone ID mismatch between parsed data and rule file"
	if err.Error() != errString {
		t.Fatalf("Expected %q but got %q", errString, err.Error())
	}
}

func TestMultipleRegions(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-error-region.yaml", &input)

	var rules RingRules
	misc.ReadYAML("../../testing/artisan-rules-1.yaml", &rules)

	_, err := rules.CalculateChanges(input, "/dev/null")
	if err == nil {
		t.Fatal("This test is expected to fail")
	}
	errString := "currently only one region is supported"
	if err.Error() != errString {
		t.Fatalf("Expected %q but got %q", errString, err.Error())
	}
}
