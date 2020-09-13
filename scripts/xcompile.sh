#!/bin/bash

# This assumes that gonative and gox is installed as per the 'one time setup' instructions
# at https://github.com/inconshreveable/gonative


CGO_ENABLED=0 gox -ldflags "${LDFLAGS}" -output="build/yq_{{.OS}}_{{.Arch}}"
# include non-default linux builds too
CGO_ENABLED=0 gox -ldflags "${LDFLAGS}" -os=linux  -output="build/yq_{{.OS}}_{{.Arch}}"

cd build
rhash -r -a . -P -o checksums

