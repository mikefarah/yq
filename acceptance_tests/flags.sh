#!/bin/bash

setUp() {
  rm test*.yml || true
  cat >test.yml <<EOL
# comment
EOL
}

testStringInterpolation() {
    X=$(./yq -n '"Mike \(3 + 4)"')
    assertEquals "Mike 7" "$X"
}

testNoStringInterpolation() {
    X=$(./yq --string-interpolation=f -n '"Mike \(3 + 4)"')
    assertEquals "Mike \(3 + 4)" "$X"
}


source ./scripts/shunit2