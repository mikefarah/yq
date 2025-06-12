#!/usr/bin/env bash

test_data='
- foo: false
'

for version in 4.45.1 4.45.2 4.45.3; do
  for command in '.[] | (select(.foo) | {"foo": .foo} // {})' '.[] | (select(.foo) | {.foo} // {})'; do
    echo ${version} "${command}"
    echo -------
    echo "${test_data}" | podman run -i --rm  mikefarah/yq:${version} -o json "${command}"
    echo -------
    echo
  done
done