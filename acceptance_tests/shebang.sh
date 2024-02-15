#!/bin/bash

setUp() {
  rm test*.yq || true
  cat >test.yq <<EOL
#!./yq
.a.b
EOL
chmod +x test.yq

rm test*.yml || true
  cat >test.yml <<EOL
a: {b: apple}
EOL
}

testCanExecYqFile() {
  read -r -d '' expected << EOM
apple
EOM
   X=$(./test.yq test.yml)
  assertEquals "$expected" "$X"
}

source ./scripts/shunit2

