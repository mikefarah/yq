#!/bin/bash
set -e

if [ "$1" == "" ]; then
  echo "Please specify at a version"
  exit 1
fi

version=$1

# validate version is in the right format
echo $version | sed -r '/v4\.[0-9][0-9]\.[0-9][0-9]?$/!{q1}'

previousVersion=$(cat cmd/version.go| sed -n 's/.*Version = "\([^"]*\)"/\1/p')

echo "Updating from $previousVersion to $version"

sed -i "s/\(.*Version =\).*/\1 \"$version\"/" cmd/version.go

go build .
actualVersion=$(./yq --version)

if [ "$actualVersion" != "yq (https://github.com/mikefarah/yq/) version $version" ]; then
    echo "Failed to update version.go"
    exit 1
else
    echo "version.go updated"
fi

 version=$version ./yq -i '.version=strenv(version) | .parts.yq.source-tag=strenv(version)' snap/snapcraft.yaml

actualSnapVersion=$(./yq '.version' snap/snapcraft.yaml)

if [ "$actualSnapVersion" != "$version" ]; then
    echo "Failed to update snapcraft"
    exit 1
else
    echo "snapcraft updated"
fi

actualSnapVersion=$(./yq '.parts.yq.source-tag' snap/snapcraft.yaml)

if [ "$actualSnapVersion" != "$version" ]; then
    echo "Failed to update snapcraft"
    exit 1
else
    echo "snapcraft updated"
fi

git add cmd/version.go snap/snapcraft.yaml
git commit -m 'Bumping version'
git tag $version
git tag -f v4