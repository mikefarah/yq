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

testOutputYamlRawDefault() {
  cat >test.yml <<EOL
a: "cat"
EOL

  X=$(./yq e  '.a' test.yml)
  assertEquals "cat" "$X"

  X=$(./yq ea '.a' test.yml)
  assertEquals "cat" "$X"
}

testOutputYamlRawOff() {
  cat >test.yml <<EOL
a: "cat"
EOL

  X=$(./yq e -r=false '.a' test.yml)
  assertEquals "\"cat\"" "$X"

  X=$(./yq ea -r=false '.a' test.yml)
  assertEquals "\"cat\"" "$X"
}

testOutputJsonRaw() {
  cat >test.yml <<EOL
a: cat
EOL

  X=$(./yq e -r --output-format=json '.a' test.yml)
  assertEquals "cat" "$X"

  X=$(./yq ea -r --output-format=json '.a' test.yml)
  assertEquals "cat" "$X"
}

testOutputJsonDefault() {
  cat >test.yml <<EOL
a: cat
EOL

  X=$(./yq e --output-format=json '.a' test.yml)
  assertEquals "\"cat\"" "$X"

  X=$(./yq ea --output-format=json '.a' test.yml)
  assertEquals "\"cat\"" "$X"
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
a: {b: {c: ["cat cat"]}}
EOL

  read -r -d '' expected << EOM
a.b.c.0 = cat cat
EOM

  X=$(./yq e --output-format=props test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea --output-format=props test.yml)
  assertEquals "$expected" "$X"
}

testOutputPropertiesDontUnwrap() {
  cat >test.yml <<EOL
a: {b: {c: ["cat cat"]}}
EOL

  read -r -d '' expected << EOM
a.b.c.0 = "cat cat"
EOM

  X=$(./yq e -r=false --output-format=props test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea -r=false --output-format=props test.yml)
  assertEquals "$expected" "$X"
}


testOutputPropertiesShort() {
  cat >test.yml <<EOL
a: {b: {c: ["cat cat"]}}
EOL

  read -r -d '' expected << EOM
a.b.c.0 = cat cat
EOM

  X=$(./yq e -o=p test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea -o=p test.yml)
  assertEquals "$expected" "$X"
}

testOutputCSV() {
  cat >test.yml <<EOL
- fruit: apple
  yumLevel: 5
- fruit: banana
  yumLevel: 4
EOL

  read -r -d '' expected << EOM
fruit,yumLevel
apple,5
banana,4
EOM

  X=$(./yq -o=c test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea -o=csv test.yml)
  assertEquals "$expected" "$X"
}

testOutputTSV() {
  cat >test.yml <<EOL
- fruit: apple
  yumLevel: 5
- fruit: banana
  yumLevel: 4
EOL

  read -r -d '' expected << EOM
fruit	yumLevel
apple	5
banana	4
EOM

  X=$(./yq -o=t test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea -o=tsv test.yml)
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

testOutputXmComplex() {
  cat >test.yml <<EOL
a: {b: {c: ["cat", "dog"], +f: meow}}
EOL

  read -r -d '' expected << EOM
<a>
  <b f="meow">
    <c>cat</c>
    <c>dog</c>
  </b>
</a>
EOM

  X=$(./yq e --output-format=x test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea --output-format=x test.yml)
  assertEquals "$expected" "$X"
}

source ./scripts/shunit2