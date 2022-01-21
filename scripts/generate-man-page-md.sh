#! /bin/bash
set -e

# note that this reqires pandoc to be installed.

cat ./pkg/yqlib/doc/operators/headers/Main.md > man.md
printf "\n# HOW IT WORKS\n" >> man.md
tail -n +2 how-it-works.md >> man.md

for f in ./pkg/yqlib/doc/*.md; do
  cat "$f" >> man.md
done
