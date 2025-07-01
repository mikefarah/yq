#!/bin/bash

# Currently, the `yq_nojson` feature must be enabled when using TinyGo.
tinygo build -no-debug -tags "yq_nolua yq_noini yq_notoml yq_noxml yq_nojson" .