#!/bin/bash

setUp() {
  rm test-fd-*.yml 2>/dev/null || true
}

tearDown() {
  rm test-fd-*.yml 2>/dev/null || true
}

createInputFiles() {
  local index
  for index in {1..80}; do
    ln -s examples/sample.yaml "test-fd-${index}.yml"
  done
}

testEvalClosesInputFiles() {
  createInputFiles

  (ulimit -n 32; GOGC=off ./yq eval 'select(false)' test-fd-*.yml)
  assertEquals 0 "$?"
}

testEvalAllClosesInputFiles() {
  createInputFiles

  (ulimit -n 32; GOGC=off ./yq eval-all 'select(false)' test-fd-*.yml)
  assertEquals 0 "$?"
}

source ./scripts/shunit2
