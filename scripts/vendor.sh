#!/bin/bash

set -e

govendor fetch github.com/op/go-logging
govendor fetch github.com/spf13/cobra
govendor fetch gopkg.in/yaml.v2
govendor sync
