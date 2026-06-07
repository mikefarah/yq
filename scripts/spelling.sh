#!/bin/bash

set -euo pipefail

GOPATH_TYPOS="$(go env GOPATH)/bin/typos"
TYPOS_CMD=""

if [ -f "${GOPATH_TYPOS}" ]; then
  TYPOS_CMD="${GOPATH_TYPOS}"
elif command -v typos >/dev/null 2>&1; then
  TYPOS_CMD="typos"
else
  echo "Error: typos not found in $(go env GOPATH)/bin or PATH."
  echo "Please run scripts/devtools.sh or ensure typos is installed correctly."
  exit 1
fi

git ls-files '*.go' '*.sh' '*.md' | "${TYPOS_CMD}" --file-list -
