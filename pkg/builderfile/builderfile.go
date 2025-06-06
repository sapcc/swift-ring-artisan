// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package builderfile

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"os/exec"
	"reflect"
	"sort"
	"time"

	"github.com/sapcc/go-bits/logg"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// File takes a path to a builder file. It tries to unpickle it
func File(builderFilename string) RingInfo {
	// // generate with ./unpickle.sh
	// cmd := exec.Command("python3", "-c", "'import json;import pickle;import sys;d=pickle.load(open(sys.argv[-1],\"rb\"));d[\"_dispersion_graph\"]=None;d[\"_replica2part2dev\"]=None;d[\"_last_part_moves\"]=None;print(json.dumps(d));'", builderFilename)
	// stdout, err := cmd.Output()
	// logg.Info("%s %s", cmd, string(stdout))
	// if err == exec.ErrNotFound {
	// 	logg.Fatal("Please install python3 to decode the pickle file.")
	// } else if err != nil {
	// 	logg.Fatal(err.Error())
	// }

	// var pickleData PickleData
	// decoder := json.NewDecoder(bytes.NewReader(stdout))
	// decoder.DisallowUnknownFields()
	// err = decoder.Decode(&pickleData)
	// if err != nil {
	// 	logg.Fatal(err.Error())
	// }

	pickleData := decodeBuilderFile(builderFilename)
	ring := RingInfo{
		ID:                    pickleData.ID,
		Version:               pickleData.Version,
		Devices:               pickleData.Devices,
		DeviceCount:           uint64(len(pickleData.Devices)),
		Dispersion:            pickleData.Dispersion,
		Partitions:            pickleData.Partitions,
		Regions:               pickleData.Devices[0].Region, // FIXME: make multi region aware
		Replicas:              pickleData.Replicas,
		OverloadFactorDecimal: pickleData.Overload,
	}
	// round to two decimal places to match the cli output
	ring.Dispersion = math.Round(ring.Dispersion*100) / 100

	// optional compare pickle parser with cli parser
	command := exec.Command("swift-ring-builder", builderFilename)
	stdout, err := command.Output()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			logg.Debug("Did not find swift-ring-builder in PATH, skipping consistency check")
			return ring
		}
		logg.Fatal("while running swift-ring-builder: " + err.Error())
	}

	ringParsed := Input(bytes.NewReader(stdout))
	// overwrite some data that the parser method but not the pickler method extracts
	ringParsed.Balance = 0
	ringParsed.FileName = ""
	ringParsed.ReassignedCooldown = 0
	ringParsed.ReassignedRemaining = time.Time{}
	ringParsed.Zones = 0
	ringParsed.OverloadFactorPercent = 0 // rely on OverloadFactorDecimal

	sort.Slice(ringParsed.Devices, func(i, j int) bool {
		return ringParsed.Devices[i].ID < ringParsed.Devices[j].ID
	})

	equal := reflect.DeepEqual(ringParsed, ring)
	if !equal {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(fmt.Sprintf("%+v\n", ringParsed), fmt.Sprintf("%+v\n", ring), false)
		logg.Info("Pickle parsed output and swift-ring-builder output are not equal. What is going on here?!")
		logg.Fatal(dmp.DiffPrettyText(diffs))
	}

	return ring
}
