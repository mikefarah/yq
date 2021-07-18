#!/bin/bash

setUp() {
  cat >test.yml <<EOL
---
a: test
EOL
}

testLeadingSeperatorPipeIntoEvalSeq() {
  X=$(cat test.yml | ./yq e -)
  expected=$(cat test.yml)
  assertEquals "$expected" "$X"
}


testLeadingSeperatorEvalSeq() {
  X=$(./yq e test.yml)
  expected=$(cat test.yml)
  assertEquals "$expected" "$X"
}

testLeadingSeperatorPipeIntoEvalAll() {
  X=$(cat test.yml | ./yq ea -)
  expected=$(cat test.yml)
  assertEquals "$expected" "$X"
}


testLeadingSeperatorEvalAll() {
  X=$(./yq ea test.yml)
  expected=$(cat test.yml)
  assertEquals "$expected" "$X"
}

testLeadingSeperatorMultiDocEval() {
  read -r -d '' expected << EOM
---
a: test
---
version: 3
application: MyApp
EOM


  X=$(./yq e '.' test.yml examples/order.yaml)
  assertEquals "$expected" "$X"
}

testLeadingSeperatorMultiDocEvalAll() {
  read -r -d '' expected << EOM
---
a: test
---
version: 3
application: MyApp
EOM


  X=$(./yq ea '.' test.yml examples/order.yaml)
  assertEquals "$expected" "$X"
}

source ./scripts/shunit2