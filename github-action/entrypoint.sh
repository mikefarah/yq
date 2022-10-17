#!/bin/sh -l
set -e
echo "::debug::\$cmd: $1"
RESULT=$(eval "$1")
RESULT="${RESULT//'%'/'%25'}"
RESULT="${RESULT//$'\n'/'%0A'}"
RESULT="${RESULT//$'\r'/'%0D'}"
echo "::debug::\$RESULT: $RESULT"
# updating from 
# https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#setting-an-output-parameter
echo "RESULT=$RESULT" >> $GITHUB_OUTPUT