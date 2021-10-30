---
description: yq is a lightweight and portable command-line YAML processor
---

# yq

&#x20;![Build](https://github.com/mikefarah/yq/workflows/Build/badge.svg) ![Docker Pulls](https://img.shields.io/docker/pulls/mikefarah/yq.svg) ![Github Releases (by Release)](https://img.shields.io/github/downloads/mikefarah/yq/total.svg) ![Go Report](https://goreportcard.com/badge/github.com/mikefarah/yq)

`yq` is a lightweight and portable command-line YAML processor. `yq` uses [jq](https://github.com/stedolan/jq) like syntax but works with yaml files as well as json. It doesn't yet support everything `jq` does - but it does support the most common operations and functions, and more is being added continuously.

`yq` is written in go - so you can download a dependency free binary for your platform and you are good to go! If you prefer there are a variety of package managers that can be used as well as docker, all listed below.

## v3 users:

Version 4 of `yq` is a major upgrade that now fully supports complex expressions for powerful filtering, slicing and dicing yaml documents. This new version uses syntax very similar to `jq` and works very similarly, so if you've used `jq` before this should be a low learning curve - however it is _quite_ different from previous versions of `yq`. Take a look at the [upgrade guide](upgrading-from-v3.md) for help.

Support for v3 will cease August 2021, until then, critical bug and security fixes will still get applied if required.

## How it works

`yq` works by running`yaml` nodes against a filter expression. The filter expression is made of of operators that pipe into each other. yaml nodes are piped through operators, operators may either return a different set of nodes (e.g. children) or modify the nodes (e.g. update values). See the [operator documentation](https://mikefarah.gitbook.io/yq/operators) for more details and examples.

```
yq eval 'select(.a == "frog") | .b.c = "dragon"' file.yaml
```

## Install

#### [Download the latest binary](https://github.com/mikefarah/yq/releases/latest)

### wget

Use wget to download the pre-compiled binaries:

#### Compressed via tar.gz

```bash
wget https://github.com/mikefarah/yq/releases/download/${VERSION}/${BINARY}.tar.gz -O - |\
  tar xz && mv ${BINARY} /usr/bin/yq
```

#### Plain binary

```bash
wget https://github.com/mikefarah/yq/releases/download/${VERSION}/${BINARY} -O /usr/bin/yq &&\
    chmod +x /usr/bin/yq
```

For instance, VERSION=v4.2.0 and BINARY=yq\_linux\_amd64

### Homebrew

[Homebrew](https://brew.sh) is a package manger for MacOS and Linux.

```
brew install yq
```

for the (deprecated) v3 version

```
brew install yq@3
```

Note that for v3, as it is a versioned brew it will not add the `yq` command to your path automatically. Please follow the instructions given by brew upon installation.

### Snap

Snap can be used on all major Linux distributions.

```
snap install yq
```

or, for the (deprecated) v3 version:

```
snap install yq --channel=v3/stable
```

**Snap notes**

`yq` installs with [_strict confinement_](https://docs.snapcraft.io/snap-confinement/6233) in snap, this means it doesn't have direct access to root files. To read root files you can:

```
sudo cat /etc/myfile | yq e '.a.path' -
```

And to write to a root file you can either use [sponge](https://linux.die.net/man/1/sponge):

```
sudo cat /etc/myfile | yq e '.a.path = "value"' - | sudo sponge /etc/myfile
```

or write to a temporary file:

```
sudo cat /etc/myfile | yq e '.a.path = "value"' | sudo tee /etc/myfile.tmp
sudo mv /etc/myfile.tmp /etc/myfile
rm /etc/myfile.tmp
```

### Docker

**Oneshot use:**

```bash
docker run --rm -v "${PWD}":/workdir mikefarah/yq <command> [flags] [expression ]FILE...
```

**Run commands interactively:**

```bash
docker run --rm -it -v "${PWD}":/workdir --entrypoint sh mikefarah/yq
```

It can be useful to have a bash function to avoid typing the whole docker command:

```bash
yq() {
  docker run --rm -i -v "${PWD}":/workdir mikefarah/yq "$@"
}
```

### Go Get

```
GO111MODULE=on go get github.com/mikefarah/yq
```

## Community Supported Installation methods

As these are supported by the community :heart: - however, they may be out of date with the officially supported releases.

## Webi

```
webi yq
```

See [webi](https://webinstall.dev)  for more infor.

Supported by @adithyasunil26 ([https://github.com/webinstall/webi-installers/tree/master/yq](https://github.com/webinstall/webi-installers/tree/master/yq))

### Chocolatey for Windows:

```
choco install yq
```

Supported by @chillum ([https://chocolatey.org/packages/yq](https://chocolatey.org/packages/yq))

### MacPorts:

```
sudo port selfupdate
sudo port install yq
```

Supported by @herbygillot ([https://ports.macports.org/maintainer/github/herbygillot](https://ports.macports.org/maintainer/github/herbygillot))

### Alpine Linux

* Enable edge/community repo by adding `$MIRROR/alpine/edge/community` to `/etc/apk/repositories`
* Update database index with `apk update`
* Install yq with `apk add yq`

Supported by Tuan Hoang [https://pkgs.alpinelinux.org/package/edge/community/x86/yq](https://pkgs.alpinelinux.org/package/edge/community/x86/yq)

#### On Ubuntu 16.04 or higher from Debian package:

```bash
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys CC86BB64
sudo add-apt-repository ppa:rmescandon/yq
sudo apt update
sudo apt install yq -y
```

Supported by @rmescandon ([https://launchpad.net/\~rmescandon/+archive/ubuntu/yq](https://launchpad.net/\~rmescandon/+archive/ubuntu/yq))

## Parsing engine and YAML spec support

Under the hood, yq uses [go-yaml v3](https://github.com/go-yaml/yaml/tree/v3) as the yaml parser, which supports [yaml spec 1.2](https://yaml.org/spec/1.2/spec.html). In particular, note that in 1.2 the values 'yes'/'no' are no longer interpreted as booleans, but as strings.
