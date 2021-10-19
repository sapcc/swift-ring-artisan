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
	"io/ioutil"
	"testing"

	"github.com/sapcc/go-bits/assert"
	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/swift-ring-artisan/pkg/parse"
	"github.com/sapcc/swift-ring-artisan/pkg/rules"
	"gopkg.in/yaml.v2"
)

func TestParse(t *testing.T) {
	inputFile, err := ioutil.ReadFile("../../testing/builder-output.yaml")
	if err != nil {
		logg.Fatal(err.Error())
	}
	var input parse.MetaData
	yaml.Unmarshal(inputFile, &input)

	expectedFile, err := ioutil.ReadFile("../../testing/artisan-rules.yaml")
	if err != nil {
		logg.Fatal(err.Error())
	}
	var expected rules.DiskRules
	yaml.Unmarshal(expectedFile, &expected)

	metaData := Convert(input, "6TB")
	if err != nil {
		logg.Fatal(err.Error())
	}

	assert.DeepEqual(t, "parsing", metaData, expected)
}
