#!/usr/bin/env bash
set -e

exp=$1
file=$2

if [ "$2" == "" ]; then
    echo "yq"
    ./yq -oj -n "$1"
    echo "jq"
    jq -n "$1"
    
else

    echo "yq"
    ./yq -oj "$1" $2
    echo "jq"
    ./yq $2 -oj | jq "$1"
fi