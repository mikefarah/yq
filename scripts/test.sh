#!/bin/bash

go test -v $(go list ./... | grep -v -E 'examples')
