#! /bin/bash
set -e

# note that this requires pandoc to be installed.

cat ./pkg/yqlib/doc/operators/headers/Main.md > man.md
printf "\n# HOW IT WORKS\n" >> man.md
tail -n +2 how-it-works.md >> man.md

for f in ./pkg/yqlib/doc/operators/*.md; do
  cat "$f" >> man.md
done

for f in ./pkg/yqlib/doc/usage/*.md; do
  cat "$f" >> man.md
done
