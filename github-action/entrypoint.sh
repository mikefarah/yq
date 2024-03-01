#!/bin/sh -l
set -e
echo "::debug::\$cmd: $1"
RESULT=$(eval "$1")
RESULT="${RESULT//'%'/'%25'}"
RESULT="${RESULT//$'\n'/'%0A'}"
RESULT="${RESULT//$'\r'/'%0D'}"
echo "::debug::\$RESULT: $RESULT"
echo "::add-mask::$RESULT"
echo ::set-output name=result::"$RESULT"
