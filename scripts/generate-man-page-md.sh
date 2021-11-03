#! /bin/bash
set -e

# note that this reqires pandoc to be installed.

cat ./pkg/yqlib/doc/headers/Main.md > man.md
printf "\n# HOW IT WORKS\n" >> man.md
tail -n +2 how-it-works.md >> man.md

for f in ./pkg/yqlib/doc/*.md; do 
  docNameWithExt="${f##*/}"
  docName="${docNameWithExt%.*}"
  docNameCap=$(echo $docName | tr [a-z] [A-Z])
  printf "\n\n# ${docNameCap}\n" >> man.md
  tail -n +2 "$f" >> man.md
done
