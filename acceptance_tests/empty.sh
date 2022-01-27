#!/bin/bash

set -eo pipefail

setUp() {
  rm test*.yml || true
  cat >test.yml <<EOL
# comment
EOL
}

testEmptyEval() {
  X=$(./yq e test.yml)
  expected=$(cat test.yml)
  assertEquals 0 $?
  assertEquals "$expected" "$X"
}

testEmptyEvalNoNewLine() {
  echo -n "#comment" >test.yml
  X=$(./yq e test.yml)
  expected=$(cat test.yml)
  assertEquals 0 $?
  assertEquals "$expected" "$X"
}

testEmptyEvalNoNewLineWithExpression() {
  echo -n "# comment" >test.yml
  X=$(./yq e '.apple = "tree"' test.yml)
  read -r -d '' expected << EOM
# comment
apple: tree
EOM
  assertEquals "$expected" "$X"
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