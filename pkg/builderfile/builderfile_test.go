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
	assert.DeepEqual(t, "parsing", expected, metaData)
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
	assert.DeepEqual(t, "parsing", expected, metaData)
}
