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
	"fmt"
	"os"

	"github.com/sapcc/go-bits/logg"
)

// WriteToStdoutOrFile writes a bytes object to stdout if filename is empty or to the file
func WriteToStdoutOrFile(data []byte, filename string) {
	if filename == "" {
		fmt.Printf("%+v", string(data))
	} else {
		err := os.WriteFile(filename, data, 0644)
		if err != nil {
			logg.Fatal(fmt.Sprintf("writing data to %s failed: %s", data, err.Error()))
		}
	}
}
