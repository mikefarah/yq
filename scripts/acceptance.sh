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

# Test leading seperator logic
expected=$(cat examples/leading-seperator.yaml)

X=$(cat examples/leading-seperator.yaml | ./yq e '.' -)
if [[ $X != $expected ]]; then
  echo "Pipe into e"
  echo "Expected $expected but was $X"
  exit 1
fi

X=$(./yq e '.' examples/leading-seperator.yaml)
expected=$(cat examples/leading-seperator.yaml)
if [[ $X != $expected ]]; then
  echo "read given file e"
  echo "Expected $expected but was $X"
  exit 1
fi

X=$(cat examples/leading-seperator.yaml | ./yq ea '.' -)
if [[ $X != $expected ]]; then
  echo "Pipe into e"
  echo "Expected $expected but was $X"
  exit 1
fi

X=$(./yq ea '.' examples/leading-seperator.yaml)
expected=$(cat examples/leading-seperator.yaml)
if [[ $X != $expected ]]; then
  echo "read given file e"
  echo "Expected $expected but was $X"
  exit 1
fi

# multidoc
read -r -d '' expected << EOM
---
a: test
---
version: 3
application: MyApp
EOM

X=$(./yq e '.' examples/leading-seperator.yaml examples/order.yaml)

if [[ $X != $expected ]]; then
  echo "Multidoc with leading seperator"
  echo "Expected $expected but was $X"
  exit 1
fi

X=$(./yq ea '.' examples/leading-seperator.yaml examples/order.yaml)

if [[ $X != $expected ]]; then
  echo "Multidoc with leading seperator"
  echo "Expected $expected but was $X"
  exit 1
fi

# handle empty files
./yq e '.' examples/empty.yaml
if [[ $? != 0 ]]; then
  echo "Expected no error when processing empty file but got one"
  exit 1
fi

cat examples/empty.yaml | ./yq e '.' -
if [[ $? != 0 ]]; then
  echo "Expected no error when processing empty stdin but got one"
  exit 1
fi


echo "--success"

set -e

rm test.yml
rm test-eval-all.yml
echo "acceptance tests passed"
