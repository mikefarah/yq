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

testInputPropertiesGitHubAction() {
  cat >test.properties <<EOL
mike.things = hello
EOL

  read -r -d '' expected << EOM
mike:
  things: hello
EOM

  X=$(cat /dev/null | ./yq e -p=props test.properties)
  assertEquals "$expected" "$X"

  X=$(cat /dev/null | ./yq ea -p=props test.properties)
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


testInputXmlStrict() {
  cat >test.yml <<EOL
<?xml version="1.0"?>
<!DOCTYPE root [
<!ENTITY writer "Catherine.">
<!ENTITY copyright "(r) Great">
]>
<root>
    <item>&writer;&copyright;</item>
</root>
EOL

  X=$(./yq -p=xml --xml-strict-mode test.yml 2>&1)
  assertEquals 1 $?
  assertEquals "Error: bad file 'test.yml': XML syntax error on line 7: invalid character entity &writer;" "$X"

  X=$(./yq ea -p=xml --xml-strict-mode test.yml 2>&1)
  assertEquals "Error: bad file 'test.yml': XML syntax error on line 7: invalid character entity &writer;" "$X"
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