#!/bin/sh

my_path="$(command -v yq)"

if [ -z "$my_path" ]; then
  echo "'yq' wasn't found in your PATH, so we don't know where to put the man pages."
  echo "Please update your PATH to include yq, and run this script again."
  exit 1
fi

# ex: ~/.local/bin/yq => ~/.local/
my_prefix="$(dirname "$(dirname "$(command -v yq)")")"
mkdir -p "$my_prefix/share/man/man1/"
cp yq.1 "$my_prefix/share/man/man1/"
