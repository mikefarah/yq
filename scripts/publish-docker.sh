#!/bin/bash
set -ex
VERSION="$(git describe --tags --abbrev=0)"
docker build \
  --target production \
  --build-arg VERSION=${VERSION} \
  -t mikefarah/yq:latest \
  -t mikefarah/yq:${VERSION} \
  .