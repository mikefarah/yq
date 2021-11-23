#!/bin/sh
set -ex
go install golang.org/x/tools/cmd/goimports
go install github.com/polyfloyd/go-errorlint
wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.37.1
wget -O- -nv https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s v2.9.1
