#!/usr/bin/env bash
echo "\"python3\", \"-c\", \"'""$(command sed unpickle.py -e '/^#/d' | sed -z -e 's|\n\n\n|;|g' -e 's|\n\n|;|g' -e 's|\n|;|g' -e 's| = |=|g' -e 's|, |,|g' -e 's|\"|\\\"|g')""'\""
