#!/bin/bash

set -o errexit
set -o pipefail

GOPATH_LINT="$(go env GOPATH)/bin/golangci-lint"
BIN_LINT="./bin/golangci-lint"
LINT_CMD=""

if [ -f "$GOPATH_LINT" ]; then
    LINT_CMD="$GOPATH_LINT"
elif [ -f "$BIN_LINT" ]; then
    LINT_CMD="$BIN_LINT"
elif command -v golangci-lint &> /dev/null; then
    # Using PATH version, ensure compatibility (see TODO)
    LINT_CMD="golangci-lint"
else
    echo "Error: golangci-lint not found in $GOPATH/bin, ./bin, or PATH."
    echo "Please run scripts/devtools.sh or ensure golangci-lint is installed correctly."
    exit 1
fi

GOFLAGS="${GOFLAGS}" "$LINT_CMD" run --verbose
