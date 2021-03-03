#!/bin/bash

set -o errexit
set -o pipefail

if command -v golangci-lint &> /dev/null
then
    golangci-lint run --timeout=5m
else
  ./bin/golangci-lint run --timeout=5m
fi

