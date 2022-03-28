#!/bin/bash

gofmt -w -s .
go mod tidy
go mod vendor
