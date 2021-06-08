#!/bin/sh -l

echo "::debug::\$cmd: $1"
RESULT=$(eval "$1")
echo "::debug::\$RESULT: $RESULT"
echo ::set-output name=result::"$RESULT"
