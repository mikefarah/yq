#!/bin/bash

setUp() {
  rm test*.yml || true
  
}

testLineCountFirstLineComment() {
  cat >test.yml <<EOL
#test123 
abc: 123 
test123: 123123 
#comment 
lalilu: lalilu
EOL

  X=$(./yq '.lalilu | line' --header-preprocess=false < test.yml)
  assertEquals "5" "$X"
}

testArrayOfDocs() {
  cat >test.yml <<EOL
---
# leading comment doc 1
a: 1
---
# leading comment doc 2
a: 2
EOL

read -r -d '' expected << EOM
- # leading comment doc 1
  a: 1
- # leading comment doc 2
  a: 2
EOM

  X=$(./yq ea '[.]' --header-preprocess=false < test.yml)
  assertEquals "$expected" "$X"

}

source ./scripts/shunit2