#!/bin/bash

testWriteInPlacePipeIn() {
  result=$(./yq e -i -n '.a' 2>&1)
  assertEquals 1 $?
  assertEquals "Error: write in place flag only applicable when giving an expression and at least one file" "$result"
}

testWriteInPlacePipeInEvalall() {
  result=$(./yq ea -i -n '.a' 2>&1)
  assertEquals 1 $?
  assertEquals "Error: write in place flag only applicable when giving an expression and at least one file" "$result"
}

testWriteInPlaceWithSplit() {
  result=$(./yq e -s "cat" -i '.a = "thing"' test.yml 2>&1)
  assertEquals 1 $?
  assertEquals "Error: write in place cannot be used with split file" "$result"
}

testWriteInPlaceWithSplitEvalAll() {
  result=$(./yq ea -s "cat" -i '.a = "thing"' test.yml 2>&1)
  assertEquals 1 $?
  assertEquals "Error: write in place cannot be used with split file" "$result"
}

testNullWithFiles() {
  result=$(./yq e -n '.a = "thing"' test.yml 2>&1)
  assertEquals 1 $?
  assertEquals "Error: cannot pass files in when using null-input flag" "$result"
}

testNullWithFilesEvalAll() {
  result=$(./yq ea -n '.a = "thing"' test.yml 2>&1)
  assertEquals 1 $?
  assertEquals "Error: cannot pass files in when using null-input flag" "$result"
}



source ./scripts/shunit2