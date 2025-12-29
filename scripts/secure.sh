#!/bin/bash

set -o errexit
set -o pipefail

OPTS=(
  -exclude-dir=vendor
  -exclude-dir=.gomodcache
  -exclude-dir=.gocache
)

command -v gosec &> /dev/null && BIN=gosec || BIN=./bin/gosec
"${BIN}" "${OPTS[@]}" "${PWD}" ./...
