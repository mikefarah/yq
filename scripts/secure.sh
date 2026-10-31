#!/bin/bash

set -o errexit
set -o pipefail

OPTS=(
  -exclude-dir=vendor
  -exclude-dir=.gomodcache
  -exclude-dir=.gocache
  --exclude=G706
)

command -v gosec &> /dev/null && BIN=gosec || BIN=./bin/gosec
"${BIN}" "${OPTS[@]}" "${PWD}" ./...
