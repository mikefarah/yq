#!/bin/bash

GITHUB_TOKEN="${GITHUB_TOKEN:?missing required input \'GITHUB_TOKEN\'}"

CURRENT="$(git describe --tags --abbrev=0)"
PREVIOUS="$(git describe --tags --abbrev=0 --always "${CURRENT}"^)"
OWNER="mikefarah"
REPO="yq"

release() {
    mapfile -t logs < <(git log --pretty=oneline --abbrev-commit "${PREVIOUS}".."${CURRENT}")
    description="$(printf '%s\n' "${logs[@]}")"
    github-release release \
        --user "$OWNER" \
        --repo "$REPO" \
        --tag "$CURRENT" \
        --description "$description" ||
            github-release edit \
                --user "$OWNER" \
                --repo "$REPO" \
                --tag "$CURRENT" \
                --description "$description"
}

upload() {
    mapfile -t files < <(find ./build -mindepth 1 -maxdepth 1)
    for file in "${files[@]}"; do
        BINARY=$(basename "${file}")
        echo "--> ${BINARY}"
        github-release upload \
            --user "$OWNER" \
            --repo "$REPO" \
            --tag "$CURRENT" \
            --name "${BINARY}" \
            --file "$file"
    done
}

release
upload
