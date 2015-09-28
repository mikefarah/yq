#!/bin/bash

gofmt -w .
golint
go test

# acceptance test
X=$(go run yaml.go  sample.yaml b.c)

if [ $X != 2 ]
  then
	echo "Failed acceptance test: expected 2 but was $X"
fi
