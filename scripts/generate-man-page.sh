#! /bin/bash
set -e

# note that this reqires pandoc to be installed.

pandoc \
  --variable=title:"YQ" \
  --variable=section:"1" \
  --variable=author:"Mike Farah" \
  --variable=header:"${MAN_HEADER}" \
  --standalone --to man man.md -o yq.1