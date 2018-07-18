#!/bin/bash

# This assumes that gonative and gox is installed as per the 'one time setup' instructions
# at https://github.com/inconshreveable/gonative

gox -ldflags "${LDFLAGS}" -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}"
gox -os=linux -arch=ppc64 -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}"
gox -os=linux -arch=ppc64le -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}"

