# yq

![Build](https://github.com/mikefarah/yq/workflows/Build/badge.svg)  ![Docker Pulls](https://img.shields.io/docker/pulls/mikefarah/yq.svg) ![Github Releases (by Release)](https://img.shields.io/github/downloads/mikefarah/yq/total.svg) ![Go Report](https://goreportcard.com/badge/github.com/mikefarah/yq)


a lightweight and portable command-line YAML processor

The aim of the project is to be the [jq](https://github.com/stedolan/jq) or sed of yaml files.

## V4 released!
V4 is now officially released, it's quite different from V3 (sorry for the migration), however it is much more similar to ```jq```, using a similar expression syntax and therefore support much more complex functionality! 

If you've been using v3 and want/need to upgrade, checkout the [upgrade guide](https://mikefarah.gitbook.io/yq/v/v4.x/upgrading-from-v3).

## Install

### [Download the latest binary](https://github.com/mikefarah/yq/releases/latest)

### MacOS:
Using [Homebrew](https://brew.sh/)
```
brew install yq
```

### Ubuntu and other Linux distros supporting `snap` packages:
```
snap install yq
```

#### Snap notes
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

### wget

Use wget to download the pre-compiled binaries:

```bash
wget https://github.com/mikefarah/yq/releases/download/${VERSION}/${BINARY} -O /usr/bin/yq &&\
    chmod +x /usr/bin/yq
```

For instance, VERSION=v4.0.0 and BINARY=yq_linux_amd64

### Run with Docker

#### Oneshot use:

```bash
docker run --rm -v "${PWD}":/workdir mikefarah/yq <command> [flags] [expression ]FILE...
```

#### Run commands interactively:

```bash
docker run --rm -it -v "${PWD}":/workdir mikefarah/yq sh
```

It can be useful to have a bash function to avoid typing the whole docker command:

```bash
yq() {
  docker run --rm -i -v "${PWD}":/workdir mikefarah/yq "$@"
}
```

### Go Get:
```
GO111MODULE=on go get github.com/mikefarah/yq
```

## Community Supported Installation methods
As these are supported by the community :heart: - however, they may be out of date with the officially supported releases.


### Windows:
```
choco install yq
```
Supported by @chillum (https://chocolatey.org/packages/yq)

### Mac:
Using [MacPorts](https://www.macports.org/)
```
sudo port selfupdate
sudo port install yq
```
Supported by @herbygillot (https://ports.macports.org/maintainer/github/herbygillot)

### Alpine Linux
- Enable edge/community repo by adding ```$MIRROR/alpine/edge/community``` to ```/etc/apk/repositories```
- Update database index with ```apk update```
- Install yq with ```apk add yq```

Supported by Tuan Hoang
https://pkgs.alpinelinux.org/package/edge/community/x86/yq


### On Ubuntu 16.04 or higher from Debian package:
```sh
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys CC86BB64
sudo add-apt-repository ppa:rmescandon/yq
sudo apt update
sudo apt install yq -y
```
Supported by @rmescandon (https://launchpad.net/~rmescandon/+archive/ubuntu/yq)

## Features
- Written in portable go, so you can download a lovely dependency free binary
- Uses similar syntax as `jq` but works with YAML and JSON files
- Fully supports multi document yaml files
- Colorized yaml output
- [Deeply traverse yaml](https://mikefarah.gitbook.io/yq/v/v4.x/traverse)
- [Sort yaml by keys](https://mikefarah.gitbook.io/yq/v/v4.x/sort-keys)
- Manipulate yaml [comments](https://mikefarah.gitbook.io/yq/comment-operators), [styling](https://mikefarah.gitbook.io/yq/style), [tags](https://mikefarah.gitbook.io/yq/tag) and [anchors and aliases](https://mikefarah.gitbook.io/yq/anchor-and-alias-operators).
- [Update yaml inplace](https://mikefarah.gitbook.io/yq/v/v4.x/commands/evaluate#flags)
- [Complex expressions to select and update](https://mikefarah.gitbook.io/yq/v/v4.x/select#select-and-update-matching-values-in-map)
- Keeps yaml formatting and comments when updating (though there are issues with whitespace)
- [Convert to/from json to yaml](https://mikefarah.gitbook.io/yq/v/v4.x/usage/convert)
- [Pipe data in by using '-'](https://mikefarah.gitbook.io/yq/v/v4.x/commands/evaluate)
- [General shell completion scripts (bash/zsh/fish/powershell)](https://mikefarah.gitbook.io/yq/v/v4.x/commands/shell-completion)

## [Usage](https://mikefarah.gitbook.io/yq/)

Check out the [documentation](https://mikefarah.gitbook.io/yq/) for more detailed and advanced usage.

```
Usage:
  yq [flags]
  yq [command]

Available Commands:
  eval             Apply expression to each document in each yaml file given in sequence
  eval-all         Loads _all_ yaml documents of _all_ yaml files and runs expression once
  help             Help about any command
  shell-completion Generate completion script

Flags:
  -C, --colors        force print with colors
  -e, --exit-status   set exit status if there are no matches or null or false is returned
  -h, --help          help for yq
  -I, --indent int    sets indent level for output (default 2)
  -i, --inplace       update the yaml file inplace of first yaml file given.
  -M, --no-colors     force print with no colors
  -N, --no-doc        Don't print document separators (---)
  -n, --null-input    Don't read input, simply evaluate the expression given. Useful for creating yaml docs from scratch.
  -P, --prettyPrint    pretty print, shorthand for '... style = ""' 
  -j, --tojson        output as json. Set indent to 0 to print json in one line.
  -v, --verbose       verbose mode
  -V, --version       Print version information and quit

Use "yq [command] --help" for more information about a command.
```

Simple Example:

```bash
yq e '.a.b | length' f1.yml f2.yml 
```

## Known Issues / Missing Features
- `yq` attempts to preserve comment positions and whitespace as much as possible, but it does not handle all scenarios (see https://github.com/go-yaml/yaml/tree/v3 for details)
