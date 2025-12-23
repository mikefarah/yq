#!/bin/bash

setUp() {
  rm test*.yml 2>/dev/null || true
  rm test*.properties 2>/dev/null || true
  rm test*.csv 2>/dev/null || true
  rm test*.tsv 2>/dev/null || true
  rm test*.xml 2>/dev/null || true
  rm test*.tf 2>/dev/null || true
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

  X=$(./yq e -p=csv test.csv)
  assertEquals "$expected" "$X"

  X=$(./yq ea -p=csv test.csv)
  assertEquals "$expected" "$X"
}

testInputCSVCustomSeparator() {
  cat >test.csv <<EOL
fruit;yumLevel
apple;5
banana;4
EOL

  read -r -d '' expected << EOM
- fruit: apple
  yumLevel: 5
- fruit: banana
  yumLevel: 4
EOM

  X=$(./yq -p=csv --csv-separator ";" test.csv)
  assertEquals "$expected" "$X"

  X=$(./yq ea -p=csv --csv-separator ";" test.csv)
  assertEquals "$expected" "$X"
}

testInputCSVNoAuto() {
  cat >test.csv <<EOL
thing1
name: cat
EOL

  read -r -d '' expected << EOM
- thing1: 'name: cat'
EOM

  X=$(./yq --csv-auto-parse=f test.csv -oy)
  assertEquals "$expected" "$X"

  X=$(./yq ea --csv-auto-parse=f test.csv -oy)
  assertEquals "$expected" "$X"
}

testInputTSVNoAuto() {
  cat >test.tsv <<EOL
thing1
name: cat
EOL

  read -r -d '' expected << EOM
- thing1: 'name: cat'
EOM

  X=$(./yq --tsv-auto-parse=f test.tsv -oy)
  assertEquals "$expected" "$X"

  X=$(./yq ea --tsv-auto-parse=f test.tsv -oy)
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

  X=$(./yq -p=csv utf8.csv)
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

  X=$(./yq e -p=t test.tsv)
  assertEquals "$expected" "$X"

  X=$(./yq ea -p=t test.tsv)
  assertEquals "$expected" "$X"
}

testInputKYaml() {
  cat >test.kyaml <<'EOL'
# leading
{
  a: 1, # a line
  # head b
  b: 2,
  c: [
    # head d
    "d", # d line
  ],
}
EOL

  read -r -d '' expected <<'EOM'
# leading
a: 1 # a line
# head b
b: 2
c:
  # head d
  - d # d line
EOM

  X=$(./yq e -p=kyaml -P test.kyaml)
  assertEquals "$expected" "$X"

  X=$(./yq ea -p=kyaml -P test.kyaml)
  assertEquals "$expected" "$X"
}




testInputXml() {
  cat >test.yml <<EOL
<cat legs="4">BiBi</cat>
EOL

  read -r -d '' expected << EOM
cat:
  +content: BiBi
  +@legs: "4"
EOM

  X=$(./yq e -p=xml test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea -p=xml test.yml)
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

  X=$(./yq e -p=xml test.xml)
  assertEquals "$expected" "$X"

  X=$(./yq ea -p=xml test.xml)
  assertEquals "$expected" "$X"
}

testInputXmlRoundtrip() {
  cat >test.yml <<EOL
<?xml version="1.0"?>
<!DOCTYPE config SYSTEM "/etc/iwatch/iwatch.dtd" >
<map xmlns="some-namespace" xmlns:xsi="some-instance" xsi:schemaLocation="some-url">Meow</map>
EOL

  read -r -d '' expected << EOM
<?xml version="1.0"?>
<!DOCTYPE config SYSTEM "/etc/iwatch/iwatch.dtd" >
<map xmlns="some-namespace" xmlns:xsi="some-instance" xsi:schemaLocation="some-url">Meow</map>
EOM

  X=$(./yq -p=xml -o=xml test.yml)
  assertEquals "$expected" "$X"

  X=$(./yq ea -p=xml -o=xml test.yml)
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

  X=$(./yq -p=xml --xml-strict-mode test.yml -o=xml 2>&1)
  assertEquals 1 $?
  assertEquals "Error: bad file 'test.yml': XML syntax error on line 7: invalid character entity &writer;" "$X"

  X=$(./yq ea -p=xml --xml-strict-mode test.yml -o=xml 2>&1)
  assertEquals "Error: bad file 'test.yml': XML syntax error on line 7: invalid character entity &writer;" "$X"
}

testInputXmlGithubAction() {
  cat >test.yml <<EOL
<cat legs="4">BiBi</cat>
EOL

  read -r -d '' expected << EOM
cat:
  +content: BiBi
  +@legs: "4"
EOM

  X=$(cat /dev/null | ./yq e -p=xml test.yml)
  assertEquals "$expected" "$X"

  X=$(cat /dev/null | ./yq ea -p=xml test.yml)
  assertEquals "$expected" "$X"
}

testInputTerraform() {
  cat >test.tf <<EOL
resource "aws_s3_bucket" "example" {
  bucket = "my-bucket"
  tags = {
    Environment = "Dev"
    Project = "Test"
  }
}
EOL

  read -r -d '' expected << EOM
resource "aws_s3_bucket" "example" {
  bucket = "my-bucket"
  tags = {
    Environment = "Dev"
    Project = "Test"
  }
}
EOM

  X=$(./yq test.tf)
  assertEquals "$expected" "$X"

  X=$(./yq ea test.tf)
  assertEquals "$expected" "$X"
}

testInputTerraformGithubAction() {
  cat >test.tf <<EOL
resource "aws_s3_bucket" "example" {
  bucket = "my-bucket"
  
  tags = {
    Environment = "Dev"
    Project = "Test"
  }
}
EOL

  read -r -d '' expected << EOM
resource "aws_s3_bucket" "example" {
  bucket = "my-bucket"
  tags = {
    Environment = "Dev"
    Project = "Test"
  }
}
EOM

  X=$(cat /dev/null | ./yq test.tf)
  assertEquals "$expected" "$X"

  X=$(cat /dev/null | ./yq ea test.tf)
  assertEquals "$expected" "$X"
}

source ./scripts/shunit2
