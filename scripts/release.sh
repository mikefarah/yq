#!/bin/bash
set -ex
GITHUB_TOKEN="${GITHUB_TOKEN:?missing required input \'GITHUB_TOKEN\'}"

CURRENT="$(git describe --tags --abbrev=0)"
PREVIOUS="$(git describe --tags --abbrev=0 --always "${CURRENT}"^)"
OWNER="mikefarah"
REPO="yq"

release() {
    github-release release \
        --user "$OWNER" \
        --draft \
        --repo "$REPO" \
        --tag "$CURRENT"
}



release

