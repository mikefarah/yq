#!/bin/bash

go test $(go list ./... | grep -v -E 'examples' | grep -v -E 'test')
