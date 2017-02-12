#!/bin/bash

gofmt -w .

./ci.sh

go install
