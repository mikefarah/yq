#!/bin/bash

setUp() {
  rm test*.yml || true
}

testBasicSplitWithName() {
  cat >test.yml <<EOL
a: test_doc1
--- 
a: test_doc2
EOL

  ./yq e test.yml -s ".a"

  doc1=$(cat test_doc1.yml)
  
  assertEquals "a: test_doc1" "$doc1"

  doc2=$(cat test_doc2.yml)
  read -r -d '' expectedDoc2 << EOM
---
a: test_doc2
EOM
  assertEquals "$expectedDoc2" "$doc2"
}

testBasicSplitWithNameEvalAll() {
  cat >test.yml <<EOL
a: test_doc1
--- 
a: test_doc2
EOL

  ./yq ea test.yml -s ".a"

  doc1=$(cat test_doc1.yml)
  
  assertEquals "a: test_doc1" "$doc1"

  doc2=$(cat test_doc2.yml)
  read -r -d '' expectedDoc2 << EOM
---
a: test_doc2
EOM
  assertEquals "$expectedDoc2" "$doc2"
}

testBasicSplitWithIndex() {
  cat >test.yml <<EOL
a: test_doc1
--- 
a: test_doc2
EOL

  ./yq e test.yml -s '"test_" + $index'

  doc1=$(cat test_0.yml)
  
  assertEquals "a: test_doc1" "$doc1"

  doc2=$(cat test_1.yml)
  read -r -d '' expectedDoc2 << EOM
---
a: test_doc2
EOM
  assertEquals "$expectedDoc2" "$doc2"
}

testBasicSplitWithIndexEvalAll() {
  cat >test.yml <<EOL
a: test_doc1
--- 
a: test_doc2
EOL

  ./yq ea test.yml -s '"test_" + $index'

  doc1=$(cat test_0.yml)
  
  assertEquals "a: test_doc1" "$doc1"

  doc2=$(cat test_1.yml)
  read -r -d '' expectedDoc2 << EOM
---
a: test_doc2
EOM
  assertEquals "$expectedDoc2" "$doc2"
}


testArraySplitWithNameNoSeparators() {
  cat >test.yml <<EOL
- name: test_fred
  age: 35
- name: test_catherine
  age: 37
EOL

  ./yq e --no-doc -s ".name"  ".[]" test.yml 

  doc1=$(cat test_fred.yml)
  read -r -d '' expectedDoc1 << EOM
name: test_fred
age: 35
EOM

  assertEquals "$expectedDoc1" "$doc1"

  doc2=$(cat test_catherine.yml)
  read -r -d '' expectedDoc2 << EOM
name: test_catherine
age: 37
EOM
  assertEquals "$expectedDoc2" "$doc2"
}

testArraySplitWithNameNoSeparatorsEvalAll() {
  cat >test.yml <<EOL
- name: test_fred
  age: 35
- name: test_catherine
  age: 37
EOL

cat >test2.yml <<EOL
- name: test_mike
  age: 564
EOL

  ./yq ea --no-doc -s ".name"  ".[]" test.yml test2.yml

  doc1=$(cat test_fred.yml)
  read -r -d '' expectedDoc1 << EOM
name: test_fred
age: 35
EOM

  assertEquals "$expectedDoc1" "$doc1"

  doc2=$(cat test_catherine.yml)
  read -r -d '' expectedDoc2 << EOM
name: test_catherine
age: 37
EOM
  assertEquals "$expectedDoc2" "$doc2"


  doc3=$(cat test_mike.yml)
  read -r -d '' expectedDoc3 << EOM
name: test_mike
age: 564
EOM
  assertEquals "$expectedDoc3" "$doc3"
}

source ./scripts/shunit2