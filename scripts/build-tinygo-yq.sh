#!/bin/bash

# Currently, the `yq_nojson` feature must be enabled when using TinyGo.
tinygo build -no-debug -tags "yq_nolua yq_noini yq_notoml yq_noxml yq_nojson yq_nocsv yq_nobase64 yq_nouri yq_noprops yq_nosh yq_noshell yq_nohcl yq_nokyaml" .
