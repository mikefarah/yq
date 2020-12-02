# yq

![Build](https://github.com/mikefarah/yq/workflows/Build/badge.svg)  ![Docker Pulls](https://img.shields.io/docker/pulls/mikefarah/yq.svg) ![Github Releases (by Release)](https://img.shields.io/github/downloads/mikefarah/yq/total.svg) ![Go Report](https://goreportcard.com/badge/github.com/mikefarah/yq)


a lightweight and portable command-line YAML processor

The aim of the project is to be the [jq](https://github.com/stedolan/jq) or sed of yaml files.

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
sudo cat /etc/myfile | yq r - a.path
```

And to write to a root file you can either use [sponge](https://linux.die.net/man/1/sponge):
```
sudo cat /etc/myfile | yq w - a.path value | sudo sponge /etc/myfile
```
or write to a temporary file:
```
sudo cat /etc/myfile | yq w - a.path value | sudo tee /etc/myfile.tmp
sudo mv /etc/myfile.tmp /etc/myfile
rm /etc/myfile.tmp
```

### wget

Use wget to download the pre-compiled binaries:

```bash
wget https://github.com/mikefarah/yq/releases/download/${VERSION}/${BINARY} -O /usr/bin/yq &&\
    chmod +x /usr/bin/yq
```

For instance, VERSION=3.4.1 and BINARY=yq_linux_amd64


### Run with Docker

#### Oneshot use:

```bash
docker run --rm -v "${PWD}":/workdir mikefarah/yq yq [flags] <command> FILE...
```

#### Run commands interactively:

```bash
docker run --rm -it -v "${PWD}":/workdir mikefarah/yq sh
```

It can be useful to have a bash function to avoid typing the whole docker command:

```bash
yq() {
  docker run --rm -i -v "${PWD}":/workdir mikefarah/yq yq "$@"
}
```

### Go Get:
```
GO111MODULE=on go get github.com/mikefarah/yq/v3
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
- [Colorize the output](https://mikefarah.gitbook.io/yq/usage/output-format#colorize-output)
- [Deep read a yaml file with a given path expression](https://mikefarah.gitbook.io/yq/commands/read#basic)
- [List matching paths of a given path expression](https://mikefarah.gitbook.io/yq/commands/read#path-only)
- [Return the lengths of arrays/object/scalars](https://mikefarah.gitbook.io/yq/commands/read#printing-length-of-the-results)
- Update a yaml file given a [path expression](https://mikefarah.gitbook.io/yq/commands/write-update#basic) or [script file](https://mikefarah.gitbook.io/yq/commands/write-update#basic)
- Update creates any missing entries in the path on the fly
- Deeply [compare](https://mikefarah.gitbook.io/yq/commands/compare) yaml files
- Keeps yaml formatting and comments when updating
- [Validate a yaml file](https://mikefarah.gitbook.io/yq/commands/validate)
- Create a yaml file given a [deep path and value](https://mikefarah.gitbook.io/yq/commands/create#creating-a-simple-yaml-file) or a [script file](https://mikefarah.gitbook.io/yq/commands/create#creating-using-a-create-script)
- [Prefix a path to a yaml file](https://mikefarah.gitbook.io/yq/commands/prefix)
- [Convert to/from json to yaml](https://mikefarah.gitbook.io/yq/usage/convert)
- [Pipe data in by using '-'](https://mikefarah.gitbook.io/yq/commands/read#from-stdin)
- [Merge](https://mikefarah.gitbook.io/yq/commands/merge) multiple yaml files with various options for [overriding](https://mikefarah.gitbook.io/yq/commands/merge#overwrite-values) and [appending](https://mikefarah.gitbook.io/yq/commands/merge#append-values-with-arrays)
- Supports multiple documents in a single yaml file for [reading](https://mikefarah.gitbook.io/yq/commands/read#multiple-documents), [writing](https://mikefarah.gitbook.io/yq/commands/write-update#multiple-documents) and [merging](https://mikefarah.gitbook.io/yq/commands/merge#multiple-documents)
- General shell completion scripts (bash/zsh/fish/powershell) (https://mikefarah.gitbook.io/yq/commands/shell-completion)

## [Usage](https://mikefarah.gitbook.io/yq/)

Check out the [documentation](https://mikefarah.gitbook.io/yq/) for more detailed and advanced usage.

```
Usage:
  yq [flags]
  yq [command]

Available Commands:
  compare          yq x [--prettyPrint/-P] dataA.yaml dataB.yaml 'b.e(name==fr*).value'
  delete           yq d [--inplace/-i] [--doc/-d index] sample.yaml 'b.e(name==fred)'
  help             Help about any command
  merge            yq m [--inplace/-i] [--doc/-d index] [--overwrite/-x] [--append/-a] sample.yaml sample2.yaml
  new              yq n [--script/-s script_file] a.b.c newValue
  prefix           yq p [--inplace/-i] [--doc/-d index] sample.yaml a.b.c
  read             yq r [--printMode/-p pv] sample.yaml 'b.e(name==fr*).value'
  shell-completion Generates shell completion scripts
  validate         yq v sample.yaml
  write            yq w [--inplace/-i] [--script/-s script_file] [--doc/-d index] sample.yaml 'b.e(name==fr*).value' newValue

Flags:
  -C, --colors        print with colors
  -h, --help          help for yq
  -I, --indent int    sets indent level for output (default 2)
  -P, --prettyPrint   pretty print
  -j, --tojson        output as json. By default it prints a json document in one line, use the prettyPrint flag to print a formatted doc.
  -v, --verbose       verbose mode
  -V, --version       Print version information and quit

Use "yq [command] --help" for more information about a command.
```

## Upgrade from V2
If you've been using v2 and want/need to upgrade, checkout the [upgrade guide](https://mikefarah.gitbook.io/yq/upgrading-from-v2).

## V4 is in development!
If you're keen - check out the alpha release [here](https://github.com/mikefarah/yq/releases/), or in docker `mikefarah/yq:4-beta1` and the docs (also in beta) [here](https://mikefarah.gitbook.io/yq/v/v4.x-alpha/).

V4 is quite different from V3 (sorry for the migration), however it will be much more similar to ```jq```, use a similar expression syntax and therefore support much more complex functionality! 

For now - new features will be held off from V3 in anticipation of the V4 build. Critical fixes / issues will still be released.

## Known Issues / Missing Features
- `yq` attempts to preserve comment positions and whitespace as much as possible, but it does not handle all scenarios (see https://github.com/go-yaml/yaml/tree/v3 for details)
- You cannot (yet) select multiple paths/keys from the yaml to be printed out (https://github.com/mikefarah/yq/issues/287) (although you can in v4!)
