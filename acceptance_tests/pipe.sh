#!/bin/bash

setUp() {
  rm test*.yml || true
  cat >test.yml <<EOL
a: frog
EOL
}

testPipeViaCatWithParam() {
  X=$(cat test.yml | ./yq '.a')
  assertEquals "frog" "$X"
}

testPipeViaCatWithParamEval() {
  X=$(cat test.yml | ./yq e '.a')
  assertEquals "frog" "$X"
}

testPipeViaCatWithParamEvalAll() {
  X=$(cat test.yml | ./yq ea '.a')
  assertEquals "frog" "$X"
}

testPipeViaCatNoParam() {
  X=$(cat test.yml | ./yq)
  assertEquals "a: frog" "$X"
}

testPipeViaCatNoParamEval() {
  X=$(cat test.yml | ./yq e)
  assertEquals "a: frog" "$X"
}

testPipeViaCatNoParamEvalAll() {
  X=$(cat test.yml | ./yq ea)
  assertEquals "a: frog" "$X"
}

testPipeViaFileishWithParam() {
  X=$(./yq '.a' < test.yml)
  assertEquals "frog" "$X"
}

testPipeViaFileishWithParamEval() {
  X=$(./yq e '.a' < test.yml)
  assertEquals "frog" "$X"
}

testPipeViaFileishWithParamEvalAll() {
  X=$(./yq ea '.a' < test.yml)
  assertEquals "frog" "$X"
}

testPipeViaFileishNoParam() {
  X=$(./yq < test.yml)
  assertEquals "a: frog" "$X"
}

testPipeViaFileishNoParamEval() {
  X=$(./yq e < test.yml)
  assertEquals "a: frog" "$X"
}

testPipeViaFileishNoParamEvalAll() {
  X=$(./yq ea < test.yml)
  assertEquals "a: frog" "$X"
}

source ./scripts/shunit2