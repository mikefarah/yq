#!/bin/bash

set -e

# acceptance test
X=$(./bin/yaml w ./examples/sample.yaml b.c 3 | ./bin/yaml r - b.c)

if [ $X != 3 ]
  then
  echo "Failed acceptance test: expected 2 but was $X"
  exit 1
fi
