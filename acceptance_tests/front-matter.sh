#!/bin/bash

setUp() {
  rm test*.yml || true
  cat >test.yml <<EOL
---
a: apple
b: cat
---
not yaml
c: at
EOL
}

testFrontMatterProcessEval() {
  read -r -d '' expected << EOM
---
a: apple
b: dog
---
not yaml
c: at
EOM
  ./yq e --front-matter="process" '.b = "dog"' test.yml -i
  assertEquals "$expected" "$(cat test.yml)"
}

testFrontMatterProcessEvalAll() {
  read -r -d '' expected << EOM
---
a: apple
b: dog
---
not yaml
c: at
EOM
  ./yq ea --front-matter="process" '.b = "dog"' test.yml -i
  assertEquals "$expected" "$(cat test.yml)"
}

testFrontMatterExtractEval() {
    cat >test.yml <<EOL
a: apple
b: cat
---
not yaml
c: at
EOL

  read -r -d '' expected << EOM
a: apple
b: dog
EOM
  ./yq e --front-matter="extract" '.b = "dog"' test.yml -i
  assertEquals "$expected" "$(cat test.yml)"
}

testFrontMatterExtractEvalAll() {
  cat >test.yml <<EOL
a: apple
b: cat
---
not yaml
c: at
EOL

  read -r -d '' expected << EOM
a: apple
b: dog
EOM
  ./yq ea --front-matter="extract" '.b = "dog"' test.yml -i
  assertEquals "$expected" "$(cat test.yml)"
}


source ./scripts/shunit2