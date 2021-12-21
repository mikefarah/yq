#!/bin/bash

setUp() {
  rm test*.yml || true
}

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

  X=$(./yq e --output-format=json test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea --output-format=json test.yml)
  assertEquals "$expected" "$X"
}

testOutputJsonShort() {
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

  X=$(./yq e -o=j test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea -o=j test.yml)
  assertEquals "$expected" "$X"
}

testOutputProperties() {
  cat >test.yml <<EOL
a: {b: {c: ["cat"]}}
EOL

  read -r -d '' expected << EOM
a.b.c.0 = cat
EOM

  X=$(./yq e --output-format=props test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea --output-format=props test.yml)
  assertEquals "$expected" "$X"
}

testOutputPropertiesShort() {
  cat >test.yml <<EOL
a: {b: {c: ["cat"]}}
EOL

  read -r -d '' expected << EOM
a.b.c.0 = cat
EOM

  X=$(./yq e -o=p test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea -o=p test.yml)
  assertEquals "$expected" "$X"
}

testOutputXml() {
  cat >test.yml <<EOL
a: {b: {c: ["cat"]}}
EOL

  read -r -d '' expected << EOM
<a>
  <b>
    <c>cat</c>
  </b>
</a>
EOM

  X=$(./yq e --output-format=xml test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea --output-format=xml test.yml)
  assertEquals "$expected" "$X"
}

testOutputXmlShort() {
  cat >test.yml <<EOL
a: {b: {c: ["cat"]}}
EOL

  read -r -d '' expected << EOM
<a>
  <b>
    <c>cat</c>
  </b>
</a>
EOM

  X=$(./yq e --output-format=x test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea --output-format=x test.yml)
  assertEquals "$expected" "$X"
}

source ./scripts/shunit2