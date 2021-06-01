#!/bin/bash

find . \( -path ./vendor \) -prune -o -name "*.go" -exec goimports -w {} \;
gofmt -w -s .
go mod tidy
go mod vendor
