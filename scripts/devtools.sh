#!/bin/sh
set -e
wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.21.0
go get golang.org/x/tools/cmd/goimports