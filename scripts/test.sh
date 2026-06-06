#!/bin/bash

go test $(go list ./... | grep -v -E 'examples' | grep -v -E 'test')

# Run after the main test suite: TestGoInstallCompatibility zips the module tree and
# must not run in parallel with pkg/yqlib tests that rewrite doc/*.md files.
go test -tags goinstall -run TestGoInstallCompatibility .
