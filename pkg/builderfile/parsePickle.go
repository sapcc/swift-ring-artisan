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
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/nlpodyssey/gopickle/pickle"
	"github.com/nlpodyssey/gopickle/types"
	"github.com/sapcc/go-bits/logg"
)

type pickleData struct {
	ID         string
	Replicas   float64
	Dispersion float64
	Devices    []DeviceInfo `mapstructure:"devs"`
	Partitions uint64       `mapstructure:"parts"`
	Version    uint64
}

func unmarshal(input interface{}) pickleData {
	pickle := guessType(input)
	var mappedData pickleData
	mapstructure.Decode(pickle, &mappedData)
	return mappedData
}

func guessType(input interface{}) interface{} {
	switch v := input.(type) {
	case *types.Dict:
		data := make(map[string]interface{})
		for _, entry := range *v {
			key := entry.Key.(string)
			// skip keys which have have tuple indexed Dicts
			if key == "_dispersion_graph" || key == "_replica2part2dev" || key == "_last_part_moves" {
				continue
			}
			// skip balance to avoid rounding errors when comparing with text based parser
			if key == "balance" {
				continue
			}
			data[entry.Key.(string)] = guessType(entry.Value)
		}
		return data
	case *types.List:
		var data []interface{}
		for _, entry := range *v {
			// skip empty entries in Devices so that mapstructure does not convert them to empty DeviceInfos
			if entry == nil {
				continue
			}
			data = append(data, guessType(entry))
		}
		return data
	case bool, float64, int, nil, string:
		return v
	default:
		logg.Fatal("Can't translate type %T", v)
	}

	return nil
}

type Array types.List

var _ types.Callable = &Array{}

func (*Array) Call(args ...interface{}) (interface{}, error) {
	// TODO: make fully functional
	// args[0] contains a type like B or H which is not taken into account
	return args[1].(*types.List), nil
}

func decodeBuilderFile(builderFilename string) pickleData {
	builderReader, err := os.Open(builderFilename)
	if err != nil {
		logg.Fatal("Reading file failed: %s", err.Error())
	}
	defer builderReader.Close()
	u := pickle.NewUnpickler(builderReader)
	u.FindClass = func(module, name string) (interface{}, error) {
		if module == "array" && name == "array" {
			return &Array{}, nil
		}
		return nil, fmt.Errorf("class not found :(")
	}
	pickled, err := u.Load()
	if err != nil {
		logg.Fatal(err.Error())
	}

	return unmarshal(pickled)
}
