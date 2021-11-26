#!/bin/sh
set -ex
go mod download golang.org/x/tools@latest
wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.43.0
wget -O- -nv https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s v2.9.1
