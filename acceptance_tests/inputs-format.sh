#!/bin/bash

setUp() {
  rm test*.yml 2>/dev/null || true
  rm test*.properties 2>/dev/null || true
  rm test*.xml 2>/dev/null || true
}

testInputProperties() {
  cat >test.properties <<EOL
mike.things = hello
EOL

  read -r -d '' expected << EOM
mike:
  things: hello
EOM

  X=$(./yq e -p=props test.properties)
  assertEquals "$expected" "$X"

  X=$(./yq ea -p=props test.properties)
  assertEquals "$expected" "$X"
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