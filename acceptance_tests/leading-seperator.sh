#!/bin/bash

setUp() {
  cat >test.yml <<EOL
---
a: test
EOL
}

testLeadingSeperatorWithDoc() {
  cat >test.yml <<EOL
# hi peeps
# cool
---
a: test
---
b: cool
EOL

  read -r -d '' expected << EOM
# hi peeps
# cool
---
a: thing
---
b: cool
EOM

  X=$(./yq e '(select(di == 0) | .a) = "thing"' - < test.yml)
  assertEquals "$expected" "$X"
}

testLeadingSeperatorPipeIntoEvalSeq() {
  X=$(./yq e - < test.yml)
  expected=$(cat test.yml)
  assertEquals "$expected" "$X"
}


testLeadingSeperatorEvalSeq() {
  X=$(./yq e test.yml)
  expected=$(cat test.yml)
  assertEquals "$expected" "$X"
}

testLeadingSeperatorPipeIntoEvalAll() {
  X=$(./yq ea - < test.yml)
  expected=$(cat test.yml)
  assertEquals "$expected" "$X"
}


testLeadingSeperatorEvalAll() {
  X=$(./yq ea test.yml)
  expected=$(cat test.yml)
  assertEquals "$expected" "$X"
}

testLeadingSeperatorMultiDocEvalSimple() {
  read -r -d '' expected << EOM
---
a: test
---
version: 3
application: MyApp
EOM


  X=$(./yq e '.' test.yml examples/order.yaml)
  assertEquals "$expected" "$X"
}

testLeadingSeperatorMultiDocInOneFile() {
  cat >test.yml <<EOL
---
# hi peeps
# cool
a: test
---
b: things
EOL
  expected=$(cat test.yml)
  X=$(./yq e '.' test.yml)
  assertEquals "$expected" "$X"
}

testLeadingSeperatorMultiDocInOneFileEvalAll() {
  cat >test.yml <<EOL
---
# hi peeps
# cool
a: test
---
b: things
EOL
  expected=$(cat test.yml)
  X=$(./yq ea '.' test.yml)
  assertEquals "$expected" "$X"
}

testLeadingSeperatorMultiDocEvalComments() {
  cat >test.yml <<EOL
# hi peeps
# cool
a: test
EOL

cat >test2.yml <<EOL
# this is another doc
# great
b: sane
EOL

  read -r -d '' expected << EOM
# hi peeps
# cool
a: test
---
# this is another doc
# great
b: sane
EOM


  X=$(./yq e '.' test.yml test2.yml)
  assertEquals "$expected" "$X"
}

testLeadingSeperatorMultiDocEvalCommentsTrailingSep() {
  cat >test.yml <<EOL
# hi peeps
# cool
---
a: test
EOL

cat >test2.yml <<EOL
# this is another doc
# great
---
b: sane
EOL

  read -r -d '' expected << EOM
# hi peeps
# cool
---
a: test
---
# this is another doc
# great
---
b: sane
EOM


  X=$(./yq e '.' test.yml test2.yml)
  assertEquals "$expected" "$X"
}

testLeadingSeperatorMultiMultiDocEvalCommentsTrailingSep() {
  cat >test.yml <<EOL
# hi peeps
# cool
---
a: test
---
a1: test2
EOL

cat >test2.yml <<EOL
# this is another doc
# great
---
b: sane
---
b2: cool
EOL

  read -r -d '' expected << EOM
# hi peeps
# cool
---
a: test
---
a1: test2
---
# this is another doc
# great
---
b: sane
---
b2: cool
EOM


  X=$(./yq e '.' test.yml test2.yml)
  assertEquals "$expected" "$X"
}

testLeadingSeperatorMultiDocEvalCommentsLeadingSep() {
  cat >test.yml <<EOL
---
# hi peeps
# cool
a: test
EOL

cat >test2.yml <<EOL
---
# this is another doc
# great
b: sane
EOL

  read -r -d '' expected << EOM
---
# hi peeps
# cool
a: test
---
# this is another doc
# great
b: sane
EOM


  X=$(./yq e '.' test.yml test2.yml)
  assertEquals "$expected" "$X"
}

testLeadingSeperatorMultiDocEvalCommentsStripComments() {
  cat >test.yml <<EOL
---
# hi peeps
# cool
a: test
---
# this is another doc
# great
b: sane
EOL

  # it will be hard to remove that top level separator
  read -r -d '' expected << EOM
a: test
---
b: sane
EOM


  X=$(./yq e --header-preprocess=false '... comments=""'  test.yml)
  assertEquals "$expected" "$X"
}

testLeadingSeperatorMultiDocEvalCommentsLeadingSepNoDocFlag() {
  cat >test.yml <<EOL
---
# hi peeps
# cool
a: test
---
# this is another doc
# great
b: sane
EOL

  # it will be hard to remove that top level separator
  read -r -d '' expected << EOM
---
# hi peeps
# cool
a: test
# this is another doc
# great
b: sane
EOM


  X=$(./yq e '.' --no-doc test.yml)
  assertEquals "$expected" "$X"
}

testLeadingSeperatorMultiDocEvalJsonFlag() {
  cat >test.yml <<EOL
---
# hi peeps
# cool
a: test
EOL

cat >test2.yml <<EOL
---
# this is another doc
# great
b: sane
EOL

  read -r -d '' expected << EOM
{
  "a": "test"
}
{
  "b": "sane"
}
EOM


  X=$(./yq e '.' -j test.yml test2.yml)
  assertEquals "$expected" "$X"
}

testLeadingSeperatorMultiDocEvalAllJsonFlag() {
  cat >test.yml <<EOL
---
# hi peeps
# cool
a: test
EOL

cat >test2.yml <<EOL
---
# this is another doc
# great
b: sane
EOL

  read -r -d '' expected << EOM
{
  "a": "test"
}
{
  "b": "sane"
}
EOM


  X=$(./yq ea '.' -j test.yml test2.yml)
  assertEquals "$expected" "$X"
}

testLeadingSeperatorMultiDocEvalAll() {
  read -r -d '' expected << EOM
---
a: test
---
version: 3
application: MyApp
EOM


  X=$(./yq ea '.' test.yml examples/order.yaml)
  assertEquals "$expected" "$X"
}

source ./scripts/shunit2