#!/bin/bash

set -eo pipefail

# You may need to go install github.com/goreleaser/goreleaser/v2@latest first
GORELEASER="goreleaser build --clean"
if [ -z "$CI" ]; then
  GORELEASER+=" --snapshot"
fi

$GORELEASER

cd build

# Remove artifacts from goreleaser
rm artifacts.json config.yaml metadata.json

find . -executable -type f | xargs -I {} tar czvf {}.tar.gz {} yq.1 -C ../scripts install-man-page.sh
tar czvf yq_man_page_only.tar.gz yq.1 -C ../scripts install-man-page.sh

rm yq_windows_386.exe.tar.gz
rm yq_windows_amd64.exe.tar.gz

zip yq_windows_386.zip yq_windows_386.exe
zip yq_windows_amd64.zip yq_windows_amd64.exe

rm yq.1

rhash -r -a . -o checksums

rhash -r -a --bsd . -o checksums-bsd

rhash --list-hashes > checksums_hashes_order

cp ../scripts/extract-checksum.sh .
