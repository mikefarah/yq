#!/bin/bash

gofmt -w .
golint
go test

# acceptance test
go build
X=$(./yaml r sample.yaml b.c)

if [ $X != 2 ]
  then
	echo "Failed acceptance test: expected 2 but was $X"
  exit 1
fi

go install
