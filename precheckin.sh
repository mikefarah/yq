#!/bin/bash

gofmt -w .
golint
go test
go build

# acceptance test
X=$(./yaml w sample.yaml b.c 3 | ./yaml r - b.c)

if [ $X != 3 ]
  then
	echo "Failed acceptance test: expected 2 but was $X"
  exit 1
fi

go install
