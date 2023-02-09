#!/bin/sh
set -ex
go mod download golang.org/x/tools@latest
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.51.1
wget -O- -nv https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s v2.15.0
