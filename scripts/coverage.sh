#!/bin/bash

set -e

go test -coverprofile=coverage.out
go tool cover -html=coverage.out -o cover/coverage.html
