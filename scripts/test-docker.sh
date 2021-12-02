#! /bin/bash
set -e

docker build . -t temp
docker run --rm -it --entrypoint sh temp -c 'touch a'
