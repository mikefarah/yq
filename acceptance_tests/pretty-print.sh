#!/bin/bash

setUp() {
  rm test*.yml || true
}

testPrettyPrintWithBooleans() {
  cat >test.yml <<EOL
leaveUnquoted: [yes, no, on, off, y, n, true, false]
leaveQuoted: ["yes", "no", "on", "off", "y", "n", "true", "false"]

EOL

  read -r -d '' expected << EOM
leaveUnquoted:
  - yes
  - no
  - on
  - off
  - y
  - n
  - true
  - false
leaveQuoted:
  - "yes"
  - "no"
  - "on"
  - "off"
  - "y"
  - "n"
  - "true"
  - "false"
EOM

  X=$(./yq e --prettyPrint test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea --prettyPrint test.yml)
  assertEquals "$expected" "$X"
}

testPrettyPrintWithBooleansCapitals() {
  cat >test.yml <<EOL
leaveUnquoted: [YES, NO, ON, OFF, Y, N, TRUE, FALSE]
leaveQuoted: ["YES", "NO", "ON", "OFF", "Y", "N", "TRUE", "FALSE"]
EOL

  read -r -d '' expected << EOM
leaveUnquoted:
  - YES
  - NO
  - ON
  - OFF
  - Y
  - N
  - TRUE
  - FALSE
leaveQuoted:
  - "YES"
  - "NO"
  - "ON"
  - "OFF"
  - "Y"
  - "N"
  - "TRUE"
  - "FALSE"
EOM

  X=$(./yq e --prettyPrint test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea --prettyPrint test.yml)
  assertEquals "$expected" "$X"
}

testPrettyPrintOtherStringValues() {
  cat >test.yml <<EOL
leaveUnquoted: [yesSir, hellno, bonapite]
makeUnquoted: ["yesSir", "hellno", "bonapite"]
EOL

  read -r -d '' expected << EOM
leaveUnquoted:
  - yesSir
  - hellno
  - bonapite
makeUnquoted:
  - yesSir
  - hellno
  - bonapite
EOM

  X=$(./yq e --prettyPrint test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea --prettyPrint test.yml)
  assertEquals "$expected" "$X"
}

testPrettyPrintKeys() {
  cat >test.yml <<EOL
"removeQuotes": "please"
EOL

  read -r -d '' expected << EOM
removeQuotes: please
EOM

  X=$(./yq e --prettyPrint test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea --prettyPrint test.yml)
  assertEquals "$expected" "$X"
}

testPrettyPrintOtherStringValues() {
  cat >test.yml <<EOL
leaveUnquoted: [yesSir, hellno, bonapite]
makeUnquoted: ["yesSir", "hellno", "bonapite"]
EOL

  read -r -d '' expected << EOM
leaveUnquoted:
  - yesSir
  - hellno
  - bonapite
makeUnquoted:
  - yesSir
  - hellno
  - bonapite
EOM

  X=$(./yq e --prettyPrint test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea --prettyPrint test.yml)
  assertEquals "$expected" "$X"
}

testPrettyPrintStringBlocks() {
  cat >test.yml <<EOL
"removeQuotes": | 
  "please"
EOL

  read -r -d '' expected << EOM
removeQuotes: |
  "please"
EOM

  X=$(./yq e --prettyPrint test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea --prettyPrint test.yml)
  assertEquals "$expected" "$X"
}

testPrettyPrintWithExpression() {
  cat >test.yml <<EOL
a: {b: {c: ["cat"]}}
EOL

  read -r -d '' expected << EOM
b:
  c:
    - cat
EOM

  X=$(./yq e '.a' --prettyPrint test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea '.a' --prettyPrint test.yml)
  assertEquals "$expected" "$X"
}

source ./scripts/shunit2