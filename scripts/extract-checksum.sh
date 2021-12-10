#!/bin/bash

# This script works with checksums_hashes_order and checksums to extract the relevant
# sha of the various yq downloads. You can then use your favourite checksum tool to validate.
# <CHECKSUM> must match an entry in checksums_hashes_order.
#
# Usage: ./extract-checksum.sh <CHECKSUM> <FILE>
# E.g: ./extract-checksum.sh SHA-256 yq_linux_amd64.tar.gz
# Outputs: 
# yq_linux_amd64.tar.gz	acebc9d07aa2d0e482969b2c080ee306e8f58efbd6f2d857eefbce6469da1473
#
# Usage with rhash:
# ./extract-checksum.sh SHA-256 yq_linux_amd64.tar.gz | rhash -c -
#
# Tip, if you want the checksum first then the filename  (e.g. for the md5sum command)
# then you can pipe the output of this script into awk to switch the fields around:
#
# ./extract-checksum.sh MD5 yq_linux_amd64.tar.gz | awk '{ print $2 " " $1}' | md5sum -c -
#
#

if [ "$1" == "" ]; then
  echo "Please specify at a hash algorithm from the checksum_hashes_order"
  echo "Usage: $0 <HASH-ALG> <FILE>"
  exit 1
fi

if [ "$2" != "" ]; then
  # so we dont match x.tar.gz when 'x' is given
  file="$2\s"
else 
  file=""
fi

if [ ! -f "checksums_hashes_order" ]; then
  echo "This script requires checksums_hashes_order to run"
  echo "Download the file from https://github.com/mikefarah/yq/releases/ for the version of yq you are trying to validate"
  exit 1
fi

if [ ! -f "checksums" ]; then
  echo "This script requires the checksums file to run"
  echo "Download the file from https://github.com/mikefarah/yq/releases/ for the version of yq you are trying to validate"
  exit 1
fi


grepMatch=$(grep -m 1 -n "$1" checksums_hashes_order)
if [ "$grepMatch" == "" ]; then
  echo "Could not find hash algorith '$1' in checksums_hashes_order"
  exit 1
fi

set -e

lineNumber=$(echo "$grepMatch" | cut -f1 -d:)

realLineNumber="$(($lineNumber + 1))"

grep "$file" checksums | sed 's/  /\t/g' | cut -f1,$realLineNumber