#!/bin/bash

setUp() {
  cat >test.yml <<EOL
# comment
EOL
}

testEmptyEval() {
  X=$(./yq e test.yml)
  assertEquals 0 $?
}

testEmptyEvalPipe() {
  X=$(./yq e - < test.yml)
  assertEquals 0 $?
}

testEmptyCommentsWithExpressionEval() {
  read -r -d '' expected << EOM
# comment
apple: tree
EOM

  X=$(./yq e '.apple="tree"' test.yml)

  assertEquals "$expected" "$X"
}

testEmptyCommentsWithExpressionEvalAll() {
  read -r -d '' expected << EOM
# comment
apple: tree
EOM

  X=$(./yq ea '.apple="tree"' test.yml)

  assertEquals "$expected" "$X"
}

testEmptyWithExpressionEval() {
  rm test.yml
  touch test.yml
  expected="apple: tree"

  X=$(./yq e '.apple="tree"' test.yml)

  assertEquals "$expected" "$X"
}

testEmptyWithExpressionEvalAll() {
  rm test.yml
  touch test.yml
  expected="apple: tree"

  X=$(./yq ea '.apple="tree"' test.yml)

  assertEquals "$expected" "$X"
}


source ./scripts/shunit2