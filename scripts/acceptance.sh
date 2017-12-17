#!/bin/bash

set -e

# acceptance test
X=$(./yq w ./examples/sample.yaml b.c 3 | ./yq r - b.c)

if [[ $X != 3 ]]; then
  echo "Failed acceptance test: expected 2 but was $X"
  exit 1
fi
echo "acceptance tests passed"
