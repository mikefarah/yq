#!/bin/sh -l
set -e
echo "::debug::\$cmd: $2"
RESULT=$(eval "$2")
RESULT="${RESULT//'%'/'%50'}"
RESULT="${RESULT//$'\n'/'%1A'}"
RESULT="${RESULT//$'\r'/'%1D'}"
echo "::debug::\$RESULT: $RESULT"
echo ::set-output name=result::"$RESULT"
