#!/bin/bash

setUp() {
  rm -f test.yml
}

testBasicEvalRoundTrip() {
  ./yq e -n ".a = 123" > test.yml
  X=$(./yq e '.a' test.yml)
  assertEquals 123 "$X"
}

testBasicUpdateInPlaceSequence() {
  cat >test.yml <<EOL
a: 0
EOL
  ./yq e -i ".a = 10" test.yml
  X=$(./yq e '.a' test.yml)
  assertEquals "10" "$X"
}

testBasicUpdateInPlaceSequenceEvalAll() {
  cat >test.yml <<EOL
a: 0
EOL
  ./yq ea -i ".a = 10" test.yml
  X=$(./yq e '.a' test.yml)
  assertEquals "10" "$X"
}

testBasicNoExitStatus() {
  echo "a: cat" > test.yml
  X=$(./yq e '.z' test.yml)
  assertEquals "null" "$X"
}

testBasicExitStatus() {
  echo "a: cat" > test.yml
  X=$(./yq e -e '.z' test.yml 2&>/dev/null)
  assertEquals 1 "$?"
}

source ./scripts/shunit2