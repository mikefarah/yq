#!/bin/bash

# Regression tests for https://github.com/mikefarah/yq/issues/465
# `yq -i` must preserve the formatting (comments, blank lines, indentation,
# spacing before inline comments) of the untouched parts of the document.

setUp() {
  rm -f test-inplace*.yml 2>/dev/null || true
}

testInPlacePreservesCommentsBlankLinesAndCommentSpacing() {
  cat >test-inplace.yml <<EOL
# Log level for operator
logLevel: 3

# Annotations applied to operator deployment
annotations: {}

# Annotations passed to operator pod(s).
podAnnotations: {}

# Node labels for pod assignment
nodeSelector:  # +doc-gen:break
  beta.kubernetes.io/os: linux
  beta.kubernetes.io/arch: amd64
EOL

  read -r -d '' expected << EOM
# Log level for operator
logLevel: 4

# Annotations applied to operator deployment
annotations: {}

# Annotations passed to operator pod(s).
podAnnotations: {}

# Node labels for pod assignment
nodeSelector:  # +doc-gen:break
  beta.kubernetes.io/os: linux
  beta.kubernetes.io/arch: amd64
EOM

  ./yq -i '.logLevel = 4' test-inplace.yml
  X=$(cat test-inplace.yml)
  assertEquals "$expected" "$X"
}

testInPlacePreservesCustomIndentationAndQuotes() {
  cat >test-inplace-indent.yml <<EOL
# top comment

parent:
    # nested comment
    child1: 'single quoted'

    child2: "double quoted"

    list:
        - one

        - two
EOL

  read -r -d '' expected << EOM
# top comment

parent:
    # nested comment
    child1: 'changed'

    child2: "double quoted"

    list:
        - one

        - two
EOM

  ./yq -i '.parent.child1 = "changed"' test-inplace-indent.yml
  X=$(cat test-inplace-indent.yml)
  assertEquals "$expected" "$X"
}

testInPlacePreservesBlankLinesAroundLiteralBlock() {
  cat >test-inplace-literal.yml <<EOL
literal: |
  line1
  line2

after: true
EOL

  read -r -d '' expected << EOM
literal: |
  line1
  line2

after: false
EOM

  ./yq -i '.after = false' test-inplace-literal.yml
  X=$(cat test-inplace-literal.yml)
  assertEquals "$expected" "$X"
}

testInPlaceDeleteKeepsSurroundingBlankLines() {
  cat >test-inplace-del.yml <<EOL
keep: 1

# comment on gone
gone: 2

after: 3
EOL

  read -r -d '' expected << EOM
keep: 1


after: 3
EOM

  ./yq -i 'del(.gone)' test-inplace-del.yml
  X=$(cat test-inplace-del.yml)
  assertEquals "$expected" "$X"
}

testInPlaceExplicitIndentFlagWins() {
  cat >test-inplace-indent-flag.yml <<EOL
parent:
    child: 1
EOL

  ./yq -i -I 8 '.parent.child = 2' test-inplace-indent-flag.yml
  X=$(cat test-inplace-indent-flag.yml)
  assertEquals "parent:
        child: 2" "$X"
}

testStdoutOutputNotChanged() {
  # without -i, output is still normalised as before
  cat >test-inplace-stdout.yml <<EOL
# comment
a: 1

b: 2
EOL

  X=$(./yq '.a = 10' test-inplace-stdout.yml)
  assertEquals "# comment
a: 10
b: 2" "$X"
}

source ./scripts/shunit2
