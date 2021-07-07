#!/bin/bash

set -o errexit
set -o pipefail

if command -v gosec &> /dev/null
then
    gosec ${PWD} ./...
else
  ./bin/gosec ${PWD} ./...
fi