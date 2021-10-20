#! /bin/bash
set -e

# note that this reqires pandoc to be installed.

cat ./pkg/yqlib/doc/headers/Main.md > man.md
printf "\n# HOW IT WORKS\n" >> man.md
cat ./pkg/yqlib/doc/aa.md >> man.md

for f in ./pkg/yqlib/doc/*.md; do 
  docNameWithExt="${f##*/}"
  docName="${docNameWithExt%.*}"
  docNameCap=$(echo $docName | tr [a-z] [A-Z])
  if [ "$docName" != "aa" ]; then
    printf "\n\n# ${docNameCap}\n" >> man.md
    cat "$f" >> man.md
  fi
  
done
