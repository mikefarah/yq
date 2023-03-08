#!/bin/bash

setUp() {
  rm test*.yml 2>/dev/null || true
  rm test*.json 2>/dev/null || true
  rm test*.properties 2>/dev/null || true
  rm test*.csv 2>/dev/null || true
  rm test*.tsv 2>/dev/null || true
  rm test*.xml 2>/dev/null || true
}

testInputJson() {
  cat >test.json <<EOL
{ "mike" : { "things": "cool" } }
EOL

  read -r -d '' expected << EOM
mike:
  things: cool
EOM

  X=$(./yq test.json)
  assertEquals "$expected" "$X"

  X=$(./yq ea test.json)
  assertEquals "$expected" "$X"
}

testInputProperties() {
  cat >test.properties <<EOL
mike.things = hello
EOL

  read -r -d '' expected << EOM
mike:
  things: hello
EOM

  X=$(./yq e test.properties)
  assertEquals "$expected" "$X"

  X=$(./yq ea test.properties)
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

  X=$(cat /dev/null | ./yq e test.properties)
  assertEquals "$expected" "$X"

  X=$(cat /dev/null | ./yq ea test.properties)
  assertEquals "$expected" "$X"
}

testInputCSV() {
  cat >test.csv <<EOL
fruit,yumLevel
apple,5
banana,4
EOL

  read -r -d '' expected << EOM
- fruit: apple
  yumLevel: 5
- fruit: banana
  yumLevel: 4
EOM

  X=$(./yq e test.csv)
  assertEquals "$expected" "$X"

  X=$(./yq ea test.csv)
  assertEquals "$expected" "$X"
}

testInputCSVUTF8() {
  read -r -d '' expected << EOM
- id: 1
  first: john
  last: smith
- id: 1
  first: jane
  last: smith
EOM

  X=$(./yq utf8.csv)
  assertEquals "$expected" "$X"
}

testInputTSV() {
  cat >test.tsv <<EOL
fruit	yumLevel
apple	5
banana	4
EOL

  read -r -d '' expected << EOM
- fruit: apple
  yumLevel: 5
- fruit: banana
  yumLevel: 4
EOM

  X=$(./yq e test.tsv)
  assertEquals "$expected" "$X"

  X=$(./yq ea test.tsv)
  assertEquals "$expected" "$X"
}




testInputXml() {
  cat >test.xml <<EOL
<cat legs="4">BiBi</cat>
EOL

  read -r -d '' expected << EOM
cat:
  +content: BiBi
  +@legs: "4"
EOM

  X=$(./yq e test.xml)
  assertEquals "$expected" "$X"

  X=$(./yq ea test.xml)
  assertEquals "$expected" "$X"
}

testInputXmlNamespaces() {
  cat >test.xml <<EOL
<?xml version="1.0"?>
<map xmlns="some-namespace" xmlns:xsi="some-instance" xsi:schemaLocation="some-url">
</map>
EOL

  read -r -d '' expected << EOM
+p_xml: version="1.0"
map:
  +@xmlns: some-namespace
  +@xmlns:xsi: some-instance
  +@xsi:schemaLocation: some-url
EOM

  X=$(./yq e test.xml)
  assertEquals "$expected" "$X"

  X=$(./yq ea test.xml)
  assertEquals "$expected" "$X"
}

testInputXmlRoundtrip() {
  cat >test.xml <<EOL
<?xml version="1.0"?>
<!DOCTYPE config SYSTEM "/etc/iwatch/iwatch.dtd" >
<map xmlns="some-namespace" xmlns:xsi="some-instance" xsi:schemaLocation="some-url">Meow</map>
EOL

  read -r -d '' expected << EOM
<?xml version="1.0"?>
<!DOCTYPE config SYSTEM "/etc/iwatch/iwatch.dtd" >
<map xmlns="some-namespace" xmlns:xsi="some-instance" xsi:schemaLocation="some-url">Meow</map>
EOM

  X=$(./yq -o=xml test.xml)
  assertEquals "$expected" "$X"

  X=$(./yq ea -o=xml test.xml)
  assertEquals "$expected" "$X"
}


testInputXmlStrict() {
  cat >test.xml <<EOL
<?xml version="1.0"?>
<!DOCTYPE root [
<!ENTITY writer "Catherine.">
<!ENTITY copyright "(r) Great">
]>
<root>
    <item>&writer;&copyright;</item>
</root>
EOL

  X=$(./yq --xml-strict-mode test.xml -o=xml 2>&1)
  assertEquals 1 $?
  assertEquals "Error: bad file 'test.xml': XML syntax error on line 7: invalid character entity &writer;" "$X"

  X=$(./yq ea --xml-strict-mode test.xml -o=xml 2>&1)
  assertEquals "Error: bad file 'test.xml': XML syntax error on line 7: invalid character entity &writer;" "$X"
}

testInputXmlGithubAction() {
  cat >test.xml <<EOL
<cat legs="4">BiBi</cat>
EOL

  read -r -d '' expected << EOM
cat:
  +content: BiBi
  +@legs: "4"
EOM

  X=$(cat /dev/null | ./yq e test.xml)
  assertEquals "$expected" "$X"

  X=$(cat /dev/null | ./yq ea test.xml)
  assertEquals "$expected" "$X"
}

source ./scripts/shunit2
