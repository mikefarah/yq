#!/bin/bash

set -e

# acceptance test



random=$((1 + $RANDOM % 10))
./yq e -n ".a = $random" > test.yml
X=$(./yq e '.a' test.yml)

if [[ $X != $random ]]; then
  echo "Failed create: expected $random but was $X"
  exit 1
fi

echo "created yaml successfully"

update=$(($random + 1))
./yq e -i ".a = $update" test.yml

X=$(./yq e '.a' test.yml)
if [[ $X != $update ]]; then
  echo "Failed to update inplace test: expected $update but was $X"
  exit 1
fi

echo "updated in place successfully"

X=$(./yq e '.z' test.yml)
echo "no exit status success"

set +e

X=$(./yq e -e '.z' test.yml)

if [[ $? != 1 ]]; then
  echo "Expected error code 1 but was $?"
  exit 1
fi

echo "exit status success"

set -e

rm test.yml

echo "acceptance tests passed"
