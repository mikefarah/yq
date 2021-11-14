#!/bin/bash

set -o errexit
set -o pipefail

if command -v golangci-lint &> /dev/null
then
    golangci-lint run --verbose
else
  ./bin/golangci-lint run --verbose
fi

