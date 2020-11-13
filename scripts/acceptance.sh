#!/bin/bash

set -e

# acceptance test
X=$(./yq e '.b.c |= 3' ./examples/sample.yaml | ./yq e '.b.c' -)

if [[ $X != 3 ]]; then
  echo "Failed acceptance test: expected 3 but was $X"
  exit 1
fi
echo "acceptance tests passed"
