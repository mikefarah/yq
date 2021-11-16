#!/bin/bash

testLoadFileNotExist() {
  result=$(./yq e -n 'load("cat.yml")' 2>&1)
  assertEquals 1 $?
  assertEquals "Error: Failed to load cat.yml: open cat.yml: no such file or directory" "$result"
}

testLoadFileExpNotExist() {
  result=$(./yq e -n 'load(.a)' 2>&1)
  assertEquals 1 $?
  assertEquals "Error: Filename expression returned nil" "$result"
}

testStrLoadFileNotExist() {
  result=$(./yq e -n 'strload("cat.yml")' 2>&1)
  assertEquals 1 $?
  assertEquals "Error: Failed to load cat.yml: open cat.yml: no such file or directory" "$result"
}

testStrLoadFileExpNotExist() {
  result=$(./yq e -n 'strload(.a)' 2>&1)
  assertEquals 1 $?
  assertEquals "Error: Filename expression returned nil" "$result"
}

source ./scripts/shunit2