#!/bin/bash

set -e

go test -coverprofile=coverage.out -v $(go list ./... | grep -v -E 'examples' | grep -v -E 'test')
go tool cover -html=coverage.out -o coverage.html
