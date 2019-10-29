#!/bin/bash

set -o errexit
set -o pipefail

./bin/golangci-lint run

# ./bin/golangci-lint \
#   --tests \
#   --vendor \
#   --disable=aligncheck \
#   --disable=gotype \
#   --disable=goconst \
#   --disable=gocyclo \
#   --deadline=300s \
#   ./...
