#! /bin/bash
set -e

for test in acceptance_tests/*.sh; do
  echo "--------------------------------------------------------------"
  echo "$test"
  echo "--------------------------------------------------------------"
  (exec $test);
done

