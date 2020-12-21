#!/bin/bash

# This assumes that gonative and gox is installed as per the 'one time setup' instructions
# at https://github.com/inconshreveable/gonative


CGO_ENABLED=0 gox -ldflags "${LDFLAGS}" -output="build/yq_{{.OS}}_{{.Arch}}"
# include non-default linux builds too
CGO_ENABLED=0 gox -ldflags "${LDFLAGS}" -os=linux  -output="build/yq_{{.OS}}_{{.Arch}}"

cd build
rhash -r -a . -P -o checksums

rhash --list-hashes > checksums_hashes_order

gzip -k --best * 

rm checksums_hashes_order.gz
rm checksums.gz
rm yq_windows_386.exe.gz
rm yq_windows_amd64.exe.gz

zip yq_windows_386.zip yq_windows_386.exe
zip yq_windows_amd64.zip yq_windows_amd64.exe