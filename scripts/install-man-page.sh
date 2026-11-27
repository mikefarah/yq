#!/bin/sh

my_path="$(command -v yq)"

if [ -z "$my_path" ]; then
  echo "'yq' wasn't found in your PATH, so we don't know where to put the man pages."
  echo "Please update your PATH to include yq, and run this script again."
  exit 1
fi

# ex: ~/.local/bin/yq => ~/.local/
# Resolve the directory physically so that /bin, which is often a symlink to
# /usr/bin, gives /usr rather than / as the install prefix.
my_prefix="$(dirname "$(CDPATH= cd "$(dirname "$my_path")" && pwd -P)")"
mkdir -p "$my_prefix/share/man/man1/"
cp yq.1 "$my_prefix/share/man/man1/"
