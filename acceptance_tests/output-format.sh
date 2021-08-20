#!/bin/bash

testOutputJsonDeprecated() {
  cat >test.yml <<EOL
a: {b: ["cat"]}
EOL

  read -r -d '' expected << EOM
{
  "a": {
    "b": [
      "cat"
    ]
  }
}
EOM

  X=$(./yq e -j test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea -j test.yml)
  assertEquals "$expected" "$X"
}

testOutputJson() {
  cat >test.yml <<EOL
a: {b: ["cat"]}
EOL

  read -r -d '' expected << EOM
{
  "a": {
    "b": [
      "cat"
    ]
  }
}
EOM

  X=$(./yq e -o=json test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea -o=json test.yml)
  assertEquals "$expected" "$X"
}

testOutputProperties() {
  cat >test.yml <<EOL
a: {b: {c: ["cat"]}}
EOL

  read -r -d '' expected << EOM
a.b.c.0 = cat
EOM

  X=$(./yq e -o=props test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea -o=props test.yml)
  assertEquals "$expected" "$X"
}

source ./scripts/shunit2