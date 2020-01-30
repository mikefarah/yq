#!/bin/sh
set -e
wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.23.1
go get golang.org/x/tools/cmd/goimports