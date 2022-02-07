#!/bin/bash

setUp() {
  rm test*.yml || true
}

testInputXml() {
  cat >test.yml <<EOL
<cat legs="4">BiBi</cat>
EOL

  read -r -d '' expected << EOM
cat:
  +content: BiBi
  +legs: "4"
EOM

  X=$(./yq e -p=xml test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea -p=xml test.yml)
  assertEquals "$expected" "$X"
}

testInputXmlGithubAction() {
  cat >test.yml <<EOL
<cat legs="4">BiBi</cat>
EOL

  read -r -d '' expected << EOM
cat:
  +content: BiBi
  +legs: "4"
EOM

  X=$(cat /dev/null | ./yq e -p=xml test.yml)
  assertEquals "$expected" "$X"

  X=$(cat /dev/null | ./yq ea -p=xml test.yml)
  assertEquals "$expected" "$X"
}

source ./scripts/shunit2