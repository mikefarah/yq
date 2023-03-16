#!/bin/bash
go build -tags yq_notoml -tags yq_noxml -tags yq_nojson -ldflags "-s -w" .