#!/bin/bash

set -e

# acceptance test



echo "test eval-sequence"
random=$((1 + $RANDOM % 10))
./yq e -n ".a = $random" > test.yml
X=$(./yq e '.a' test.yml)

if [[ $X != $random ]]; then
  echo "Failed create: expected $random but was $X"
  exit 1
fi

echo "--success"

echo "test update-in-place"

update=$(($random + 1))
./yq e -i ".a = $update" test.yml

X=$(./yq e '.a' test.yml)
if [[ $X != $update ]]; then
  echo "Failed to update inplace test: expected $update but was $X"
  exit 1
fi

echo "--success"

echo "test eval-all"
./yq ea -n ".a = $random" > test-eval-all.yml
Y=$(./yq ea '.a' test-eval-all.yml)

if [[ $Y != $random ]]; then
  echo "Failed create with eval all: expected $random but was $X"
  exit 1
fi
echo "--success"

echo "test no exit status"
./yq e '.z' test.yml
echo "--success"

echo "test exit status"
set +e

./yq e -e '.z' test.yml

if [[ $? != 1 ]]; then
  echo "Expected error code 1 but was $?"
  exit 1
fi

echo "--success"

set -e

rm test.yml
rm test-eval-all.yml
echo "acceptance tests passed"
