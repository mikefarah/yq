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
        --repo "$REPO" \
        --tag "$CURRENT"
}

upload() {
    while IFS=  read -r -d $'\0'; do
        file=$REPLY
        BINARY=$(basename "${file}")
        echo "--> ${BINARY}"
        github-release upload \
            --replace \
            --user "$OWNER" \
            --repo "$REPO" \
            --tag "$CURRENT" \
            --name "${BINARY}" \
            --file "$file"
    done < <(find ./build -mindepth 1 -maxdepth 1 -print0)
}

release
upload
