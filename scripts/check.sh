#!/bin/bash

set -o errexit
set -o pipefail

gometalinter \
  --skip=examples \
  --tests \
  --vendor \
  --disable=aligncheck \
  --disable=gotype \
  --disable=goconst \
  --cyclo-over=20 \
  --deadline=300s \
  ./...

gometalinter \
  --skip=examples \
  --tests \
  --vendor \
  --disable=aligncheck \
  --disable=gotype \
  --disable=goconst \
  --disable=gocyclo \
  --deadline=300s \
  ./...
