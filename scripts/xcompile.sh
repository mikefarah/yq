#!/bin/bash
set -e
# you may need to go install github.com/mitchellh/gox@v1.0.1 first
echo $VERSION
CGO_ENABLED=0 gox -ldflags "-s -w ${LDFLAGS}" -output="build/yq_{{.OS}}_{{.Arch}}" --osarch="darwin/amd64 darwin/arm64 freebsd/386 freebsd/amd64 freebsd/arm linux/386 linux/amd64 linux/arm linux/arm64 linux/mips linux/mips64 linux/mips64le linux/mipsle linux/ppc64 linux/ppc64le linux/s390x netbsd/386 netbsd/amd64 netbsd/arm openbsd/386 openbsd/amd64 windows/386 windows/amd64"

cd build

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
