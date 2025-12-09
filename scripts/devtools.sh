#!/bin/sh
set -ex
go mod download golang.org/x/tools@latest
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.6.1
arch="$(uname -m)"
if [ "$arch" = "ppc64le" ]; then
	go install github.com/securego/gosec/v2/cmd/gosec@latest
else
	curl -sSfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s
fi