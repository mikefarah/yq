#!/bin/bash

set -e

gofmt -w .
golint
./ci.sh

go install
