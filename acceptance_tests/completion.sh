#!/bin/bash

testCompletionRuns() {
    result=$(./yq __complete "" 2>&1)
    assertEquals 0 $?
    assertContains "$result" "Completion ended with directive:"
}

source ./scripts/shunit2
