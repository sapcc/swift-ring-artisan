// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package builderfile

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/nlpodyssey/gopickle/pickle"
	"github.com/nlpodyssey/gopickle/types"
	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/go-bits/must"
)

type pickleData struct {
	ID         string
	Replicas   float64
	Dispersion float64
	Devices    []DeviceInfo `mapstructure:"devs"`
	Partitions uint64       `mapstructure:"parts"`
	Version    uint64
	Overload   float64
}

func unmarshal(input any) pickleData {
	var mappedData pickleData
	must.Succeed(mapstructure.Decode(guessType(input), &mappedData))
	return mappedData
}

func guessType(input any) any {
	switch v := input.(type) {
	case *types.Dict:
		data := make(map[string]any)
		for _, entry := range *v {
			key := entry.Key.(string)
			// skip keys which have tuple indexed Dicts
			if key == "_dispersion_graph" || key == "_replica2part2dev" || key == "_last_part_moves" {
				continue
			}
			// skip balance to avoid rounding errors when comparing with text based parser
			if key == "balance" {
				continue
			}
			if key == "meta" {
				var meta *map[string]string

				if value := entry.Value.(string); value != "" {
					err := json.Unmarshal([]byte(value), &meta)
					if err != nil {
						logg.Fatal("Unmarshalling meta failed: ", err.Error())
					}
				}
				data[entry.Key.(string)] = meta
				continue
			}
			data[entry.Key.(string)] = guessType(entry.Value)
		}
		return data
	case *types.List:
		var data []any
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

func (*Array) Call(args ...any) (any, error) {
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
	u.FindClass = func(module, name string) (any, error) {
		if module == "array" && name == "array" {
			return &Array{}, nil
		}
		return nil, errors.New("class not found :(")
	}
	pickled := must.Return(u.Load())

	return unmarshal(pickled)
}
