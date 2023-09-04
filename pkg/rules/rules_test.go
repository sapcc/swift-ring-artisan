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
	"testing"

	"github.com/sapcc/go-bits/assert"

	"github.com/sapcc/swift-ring-artisan/pkg/builderfile"
	"github.com/sapcc/swift-ring-artisan/pkg/misc"
)

func TestApplyRules1(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-1.yaml", &input)

	var ring RingRules
	misc.ReadYAML("../../testing/artisan-rules-changes-1.yaml", &ring)

	commandQueue, confirmations, err := ring.CalculateChanges(input, "/dev/null")
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.DeepEqual(t, "parsing", commandQueue, []string{
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-01 --weight 100 166",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-02 --weight 100 166",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-03 --weight 100 166",
	})
	assert.DeepEqual(t, "parsing", confirmations, []string(nil))
}

func TestApplyRules2(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-1.yaml", &input)

	var ring RingRules
	misc.ReadYAML("../../testing/artisan-rules-changes-2.yaml", &ring)

	commandQueue, confirmations, err := ring.CalculateChanges(input, "/dev/null")
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.DeepEqual(t, "parsing", commandQueue, []string{
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-01 --weight 100 166",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-02 --weight 100 166",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-03 --weight 100 166",
	})
	assert.DeepEqual(t, "parsing", confirmations, []string(nil))
}

func TestApplyRules3(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-2.yaml", &input)

	var ring RingRules
	misc.ReadYAML("../../testing/artisan-rules-changes-3.yaml", &ring)

	commandQueue, confirmations, err := ring.CalculateChanges(input, "/dev/null")
	if err != nil {
		t.Fatal(err.Error())
	}

	var expectedCommands []string
	for i := 1; i <= 40; i++ {
		expectedCommands = append(expectedCommands, fmt.Sprintf("swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.46.14.52 --port 6001 --device swift-%02d --weight 100 200", i))
	}
	for i := 1; i <= 40; i++ {
		expectedCommands = append(expectedCommands, fmt.Sprintf("swift-ring-builder /dev/null set_weight --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-%02d --weight 100 150", i))
	}
	assert.DeepEqual(t, "parsing", commandQueue, expectedCommands)
	assert.DeepEqual(t, "parsing", confirmations, []string(nil))
}

func TestAddDisk1(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-1.yaml", &input)

	var ring RingRules
	misc.ReadYAML("../../testing/artisan-addition-1.yaml", &ring)

	commandQueue, confirmations, err := ring.CalculateChanges(input, "/dev/null")
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.DeepEqual(t, "parsing", commandQueue, []string{
		"swift-ring-builder /dev/null add --region 1 --zone 1 --ip 10.114.1.204 --port 6001 --device swift-01 --weight 100",
		"swift-ring-builder /dev/null add --region 1 --zone 1 --ip 10.114.1.204 --port 6001 --device swift-02 --weight 100",
		"swift-ring-builder /dev/null add --region 1 --zone 1 --ip 10.114.1.204 --port 6001 --device swift-03 --weight 100",
	})
	assert.DeepEqual(t, "parsing", confirmations, []string(nil))
}

func TestAddDisk2(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-2.yaml", &input)

	var ring RingRules
	misc.ReadYAML("../../testing/artisan-addition-2.yaml", &ring)

	commandQueue, confirmations, err := ring.CalculateChanges(input, "/dev/null")
	if err != nil {
		t.Fatal(err.Error())
	}

	var expectedCommands []string
	for i := 1; i <= 12; i++ {
		expectedCommands = append(expectedCommands, fmt.Sprintf("swift-ring-builder /dev/null add --region 1 --zone 4 --ip 10.46.14.161 --port 6001 --device swift-%02d --weight 166 --meta {\"hostname\":\"node1\"}", i))
	}
	for i := 1; i <= 12; i++ {
		expectedCommands = append(expectedCommands, fmt.Sprintf("swift-ring-builder /dev/null add --region 1 --zone 4 --ip 10.46.14.248 --port 6001 --device swift-%02d --weight 166", i))
	}
	for i := 1; i <= 12; i++ {
		expectedCommands = append(expectedCommands, fmt.Sprintf("swift-ring-builder /dev/null add --region 1 --zone 4 --ip 10.46.14.42 --port 6001 --device swift-%02d --weight 166", i))
	}
	assert.DeepEqual(t, "parsing", commandQueue, expectedCommands)
	assert.DeepEqual(t, "parsing", confirmations, []string(nil))
}

func TestSetWeigthZero(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-1.yaml", &input)

	var ring RingRules
	misc.ReadYAML("../../testing/artisan-rules-zero-weight.yaml", &ring)

	commandQueue, confirmations, err := ring.CalculateChanges(input, "/dev/null")
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.DeepEqual(t, "parsing", commandQueue, []string{
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.202 --port 6001 --device swift-01 --weight 100 0",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.202 --port 6001 --device swift-02 --weight 100 0",
		"swift-ring-builder /dev/null set_weight --region 1 --zone 1 --ip 10.114.1.202 --port 6001 --device swift-03 --weight 100 0",
	})
	assert.DeepEqual(t, "parsing", confirmations, []string(nil))
}

func TestDeleteDisk1(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-1.yaml", &input)

	var ring RingRules
	misc.ReadYAML("../../testing/artisan-deletion-1.yaml", &ring)

	commandQueue, confirmations, err := ring.CalculateChanges(input, "/dev/null")
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.DeepEqual(t, "parsing", commandQueue, []string{
		"swift-ring-builder /dev/null remove --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-01 --weight 100",
		"swift-ring-builder /dev/null remove --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-02 --weight 100",
		"swift-ring-builder /dev/null remove --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-03 --weight 100",
	})
	assert.DeepEqual(t, "parsing", confirmations, []string{
		"Do you want to remove disk swift-01 on node 10.114.1.203 without first scaling its weight to 0? This poses a data loss risk.",
		"Do you want to remove disk swift-02 on node 10.114.1.203 without first scaling its weight to 0? This poses a data loss risk.",
		"Do you want to remove disk swift-03 on node 10.114.1.203 without first scaling its weight to 0? This poses a data loss risk.",
	})
}

func TestDeleteDisk2(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-3.yaml", &input)

	var ring RingRules
	misc.ReadYAML("../../testing/artisan-deletion-1.yaml", &ring)

	commandQueue, confirmations, err := ring.CalculateChanges(input, "/dev/null")
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.DeepEqual(t, "parsing", commandQueue, []string{
		"swift-ring-builder /dev/null remove --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-01 --weight 0",
		"swift-ring-builder /dev/null remove --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-02 --weight 0",
		"swift-ring-builder /dev/null remove --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-03 --weight 0",
	})
	assert.DeepEqual(t, "parsing", confirmations, []string(nil))
}

func TestDeleteDisk4(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-2.yaml", &input)

	var ring RingRules
	misc.ReadYAML("../../testing/artisan-deletion-2.yaml", &ring)

	commandQueue, confirmations, err := ring.CalculateChanges(input, "/dev/null")
	if err != nil {
		t.Fatal(err.Error())
	}

	var expectedCommands []string
	for i := 1; i <= 12; i++ {
		expectedCommands = append(expectedCommands, fmt.Sprintf("swift-ring-builder /dev/null remove --region 1 --zone 2 --ip 10.246.192.68 --port 6001 --device swift-%02d --weight 133", i))
	}
	for i := 1; i <= 12; i++ {
		expectedCommands = append(expectedCommands, fmt.Sprintf("swift-ring-builder /dev/null remove --region 1 --zone 2 --ip 10.246.192.69 --port 6001 --device swift-%02d --weight 133", i))
	}
	for i := 1; i <= 12; i++ {
		expectedCommands = append(expectedCommands, fmt.Sprintf("swift-ring-builder /dev/null remove --region 1 --zone 2 --ip 10.246.192.70 --port 6001 --device swift-%02d --weight 133", i))
	}
	for i := 1; i <= 40; i++ {
		expectedCommands = append(expectedCommands, fmt.Sprintf("swift-ring-builder /dev/null remove --region 1 --zone 2 --ip 10.46.14.116 --port 6001 --device swift-%02d --weight 100", i))
	}

	var expectedConfirmations []string
	for i := 1; i <= 12; i++ {
		expectedConfirmations = append(expectedConfirmations, fmt.Sprintf("Do you want to remove disk swift-%02d on node 10.246.192.68 without first scaling its weight to 0? This poses a data loss risk.", i))
	}
	for i := 1; i <= 12; i++ {
		expectedConfirmations = append(expectedConfirmations, fmt.Sprintf("Do you want to remove disk swift-%02d on node 10.246.192.69 without first scaling its weight to 0? This poses a data loss risk.", i))
	}
	for i := 1; i <= 12; i++ {
		expectedConfirmations = append(expectedConfirmations, fmt.Sprintf("Do you want to remove disk swift-%02d on node 10.246.192.70 without first scaling its weight to 0? This poses a data loss risk.", i))
	}
	for i := 1; i <= 40; i++ {
		expectedConfirmations = append(expectedConfirmations, fmt.Sprintf("Do you want to remove disk swift-%02d on node 10.46.14.116 without first scaling its weight to 0? This poses a data loss risk.", i))
	}
	assert.DeepEqual(t, "parsing", commandQueue, expectedCommands)
	assert.DeepEqual(t, "parsing", confirmations, expectedConfirmations)
}

func TestDeleteBrokenDisk1(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-1.yaml", &input)

	var ring RingRules
	misc.ReadYAML("../../testing/artisan-broken-1.yaml", &ring)

	commandQueue, confirmations, err := ring.CalculateChanges(input, "/dev/null")
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.DeepEqual(t, "parsing", commandQueue, []string{
		"swift-ring-builder /dev/null remove --region 1 --zone 1 --ip 10.114.1.203 --port 6001 --device swift-02 --weight 100",
	})
	assert.DeepEqual(t, "parsing", confirmations, []string(nil))
}

func TestSetOverload(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-1.yaml", &input)

	var ring RingRules
	misc.ReadYAML("../../testing/artisan-rules-overload.yaml", &ring)

	commandQueue, confirmations, err := ring.CalculateChanges(input, "/dev/null")
	if err != nil {
		t.Fatal(err.Error())
	}

	var expectedCommands []string
	expectedCommands = append(expectedCommands, fmt.Sprintf("swift-ring-builder /dev/null set_overload %.6f", 0.1))
	assert.DeepEqual(t, "parsing", commandQueue, expectedCommands)
	assert.DeepEqual(t, "parsing", confirmations, []string(nil))
}

func TestZoneMismatch(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-zone-mismatch.yaml", &input)

	var ring RingRules
	misc.ReadYAML("../../testing/artisan-rules-1.yaml", &ring)

	_, _, err := ring.CalculateChanges(input, "/dev/null")
	if err == nil {
		t.Fatal("This test is expected to fail")
	}
	errString := "zone ID mismatch between parsed data 2 and rule file 1"
	if err.Error() != errString {
		t.Fatalf("Expected %q but got %q", errString, err.Error())
	}
}

func TestMultipleRegions(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-error-region.yaml", &input)

	var ring RingRules
	misc.ReadYAML("../../testing/artisan-rules-1.yaml", &ring)

	_, _, err := ring.CalculateChanges(input, "/dev/null")
	if err == nil {
		t.Fatal("This test is expected to fail")
	}
	errString := "currently only one region is supported"
	if err.Error() != errString {
		t.Fatalf("Expected %q but got %q", errString, err.Error())
	}
}

func TestParseReplicationMismatch(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-replication-mismatch.yaml", &input)

	var ring RingRules
	misc.ReadYAML("../../testing/artisan-rules-1.yaml", &ring)

	_, _, err := ring.CalculateChanges(input, "/dev/null")
	if err == nil {
		t.Fatal("This test is expected to fail")
	}
	errString := "replication ip and port do not match with the normal ip and port which is required"
	if err.Error() != errString {
		t.Fatalf("Expected %q but got %q", errString, err.Error())
	}
}

func TestPortMismatch(t *testing.T) {
	var input builderfile.RingInfo
	misc.ReadYAML("../../testing/builder-output-1.yaml", &input)

	var ring RingRules
	misc.ReadYAML("../../testing/artisan-rules-port-mismatch.yaml", &ring)

	_, _, err := ring.CalculateChanges(input, "/dev/null")
	if err == nil {
		t.Fatal("This test is expected to fail")
	}
	errString := "port mismatch between parsed data 6001 and rule file 6002"
	if err.Error() != errString {
		t.Fatalf("Expected %q but got %q", errString, err.Error())
	}
}
