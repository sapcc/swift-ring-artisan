#!/usr/bin/env python3

# SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company
# SPDX-License-Identifier: Apache-2.0

import json
import pickle
import sys


d = pickle.load(open(sys.argv[-1], "rb"))

# Can't be trivially parsed into json and we don't need so we throw it out the window
d["_dispersion_graph"] = None
d["_replica2part2dev"] = None
d["_last_part_moves"] = None
print(json.dumps(d))
