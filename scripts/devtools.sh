#!/bin/sh
set -ex
go mod download golang.org/x/tools@v0.44.0
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/6008b81b81c690c046ffc3fd5bce896da715d5fd/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.11.3
curl -sSfL https://raw.githubusercontent.com/securego/gosec/424fc4cd9c82ea0fd6bee9cd49c2db2c3cc0c93f/install.sh | sh -s v2.22.11