#!/bin/bash
go build -tags "yq_nolua yq_noini yq_notoml yq_noxml yq_nojson yq_nohcl yq_nokyaml" -ldflags "-s -w" .
