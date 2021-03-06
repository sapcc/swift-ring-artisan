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

package misc

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sapcc/go-bits/logg"
	"gopkg.in/yaml.v2"
)

// WriteToStdoutOrFile writes a bytes object to stdout if filename is empty or to the file
func WriteToStdoutOrFile(data []byte, filename string) {
	if filename == "" {
		fmt.Printf("%+v", string(data))
	} else {
		err := os.WriteFile(filename, data, 0644)
		if err != nil {
			logg.Fatal("writing data to %s failed: %s", data, err.Error())
		}
	}
}

func ReadYAML(filename string, variable interface{}) {
	ruleFile, err := os.ReadFile(filename)
	if err != nil {
		logg.Fatal(err.Error())
	}
	err = yaml.UnmarshalStrict(ruleFile, variable)
	if err != nil {
		logg.Fatal("Parsing file %s failed: %s", filename, err.Error())
	}
}

func ParseUint(string string) uint64 {
	uint, err := strconv.ParseUint(string, 10, 64)
	if err != nil {
		logg.Fatal(err.Error())
	}
	return uint
}

func ParseFloat(string string) float64 {
	float, err := strconv.ParseFloat(string, 64)
	if err != nil {
		logg.Fatal(err.Error())
	}
	return float
}

func Contains(list []string, searchFor string) bool {
	for _, entry := range list {
		if entry == searchFor {
			return true
		}
	}
	return false
}

func AskConfirmation(text string) bool {
	fmt.Printf("%s [y/N]: ", text)
	response, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		logg.Fatal(err.Error())
	}

	response = strings.ToLower(strings.TrimSpace(response))
	if response == "y" || response == "yes" {
		return true
	} else {
		return false
	}
}
