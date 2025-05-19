#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company
# SPDX-License-Identifier: Apache-2.0

echo "\"python3\", \"-c\", \"'""$(command sed unpickle.py -e '/^#/d' | sed -z -e 's|\n\n\n|;|g' -e 's|\n\n|;|g' -e 's|\n|;|g' -e 's| = |=|g' -e 's|, |,|g' -e 's|\"|\\\"|g')""'\""
