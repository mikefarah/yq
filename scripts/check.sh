#!/bin/bash

set -o errexit
set -o pipefail

if command -v golangci-lint &> /dev/null
then
    golangci-lint run --timeout=5m
else
  ./bin/golangci-lint run --timeout=5m
fi

# ./bin/golangci-lint \
#   --tests \
#   --vendor \
#   --disable=aligncheck \
#   --disable=gotype \
#   --disable=goconst \
#   --disable=gocyclo \
#   --deadline=300s \
#   ./...
