#!/bin/bash

tearDown() {
  set -e
}

testWriteInPlacePipeIn() {
  set +e
  result=$(./yq e -i -n '.a' 2>&1)
  assertEquals 1 $?
  assertEquals "Error: write inplace flag only applicable when giving an expression and at least one file" "$result"
}

testWriteInPlacePipeInEvalall() {
  set +e
  result=$(./yq ea -i -n '.a' 2>&1)
  assertEquals 1 $?
  assertEquals "Error: write inplace flag only applicable when giving an expression and at least one file" "$result"
}

testWriteInPlaceWithSplit() {
  set +e
  result=$(./yq e -s "cat" -i '.a = "thing"' test.yml 2>&1)
  assertEquals 1 $?
  assertEquals "Error: write inplace cannot be used with split file" "$result"
}

testWriteInPlaceWithSplitEvalAll() {
  set +e
  result=$(./yq ea -s "cat" -i '.a = "thing"' test.yml 2>&1)
  assertEquals 1 $?
  assertEquals "Error: write inplace cannot be used with split file" "$result"
}

testNullWithFiles() {
  set +e
  result=$(./yq e -n '.a = "thing"' test.yml 2>&1)
  assertEquals 1 $?
  assertEquals "Error: cannot pass files in when using null-input flag" "$result"
}

testNullWithFilesEvalAll() {
  set +e
  result=$(./yq ea -n '.a = "thing"' test.yml 2>&1)
  assertEquals 1 $?
  assertEquals "Error: cannot pass files in when using null-input flag" "$result"
}



source ./scripts/shunit2