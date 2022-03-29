# yq

![Build](https://github.com/mikefarah/yq/workflows/Build/badge.svg)  ![Docker Pulls](https://img.shields.io/docker/pulls/mikefarah/yq.svg) ![Github Releases (by Release)](https://img.shields.io/github/downloads/mikefarah/yq/total.svg) ![Go Report](https://goreportcard.com/badge/github.com/mikefarah/yq)


a lightweight and portable command-line YAML, JSON and XML processor. `yq` uses [jq](https://github.com/stedolan/jq) like syntax but works with yaml files as well as json and xml. It doesn't yet support everything `jq` does - but it does support the most common operations and functions, and more is being added continuously.

yq is written in go - so you can download a dependency free binary for your platform and you are good to go! If you prefer there are a variety of package managers that can be used as well as Docker and Podman, all listed below.

## Notice for v4.x versions prior to 4.18.1
Since 4.18.1, yq's 'eval/e' command is the _default_ command and no longer needs to be specified.

Older versions will still need to specify 'eval/e'.

Similarly, '-' is no longer required as a filename to read from STDIN (unless reading from one or more files).

TLDR:

Prior to 4.18.1 
```bash
yq e '.cool' - < file.yaml
```

4.18+ 
```bash
yq '.cool' < file.yaml
```

When merging multiple files together, `eval-all/ea` is still required to tell `yq` to run the expression against all the document at once.

## Quick Usage Guide

Read a value:
```bash
yq '.a.b[0].c' file.yaml
```

Pipe from STDIN:
```bash
yq '.a.b[0].c' < file.yaml
```

Update a yaml file, inplace
```bash
yq -i '.a.b[0].c = "cool"' file.yaml
```

Update using environment variables
```bash
NAME=mike yq -i '.a.b[0].c = strenv(NAME)' file.yaml
```

Merge multiple files
```bash
# note the use of `ea` to evaluate all the files at once
# instead of in sequence
yq ea '. as $item ireduce ({}; . * $item )' path/to/*.yml
```

Multiple updates to a yaml file
```bash
yq -i '
  .a.b[0].c = "cool" |
  .x.y.z = "foobar" |
  .person.name = strenv(NAME)
' file.yaml
```

Convert JSON to YAML
```bash
yq -P sample.json
```

See the [documentation](https://mikefarah.gitbook.io/yq/) for more examples.

## Install

### [Download the latest binary](https://github.com/mikefarah/yq/releases/latest)

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

For instance, VERSION=v4.2.0 and BINARY=yq_linux_amd64

### MacOS / Linux via Homebrew:
Using [Homebrew](https://brew.sh/)
```
brew install yq
```

### Linux via snap:
```
snap install yq
```

#### Snap notes
`yq` installs with [_strict confinement_](https://docs.snapcraft.io/snap-confinement/6233) in snap, this means it doesn't have direct access to root files. To read root files you can:

```
sudo cat /etc/myfile | yq '.a.path'
```

And to write to a root file you can either use [sponge](https://linux.die.net/man/1/sponge):
```
sudo cat /etc/myfile | yq '.a.path = "value"' | sudo sponge /etc/myfile
```
or write to a temporary file:
```
sudo cat /etc/myfile | yq '.a.path = "value"' | sudo tee /etc/myfile.tmp
sudo mv /etc/myfile.tmp /etc/myfile
rm /etc/myfile.tmp
```

### Run with Docker or Podman
#### Oneshot use:

```bash
docker run --rm -v "${PWD}":/workdir mikefarah/yq [command] [flags] [expression ]FILE...
```

Note that you can run `yq` in docker without network access and other privileges if you desire,
namely `--security-opt=no-new-privileges --cap-drop all --network none`.

```bash
podman run --rm -v "${PWD}":/workdir mikefarah/yq [command] [flags] [expression ]FILE...
```

#### Pipe in via STDIN:

You'll need to pass the `-i\--interactive` flag to docker:

```bash
docker run -i --rm mikefarah/yq '.this.thing' < myfile.yml
```

```bash
podman run -i --rm mikefarah/yq '.this.thing' < myfile.yml
```

#### Run commands interactively:

```bash
docker run --rm -it -v "${PWD}":/workdir --entrypoint sh mikefarah/yq
```

```bash
podman run --rm -it -v "${PWD}":/workdir --entrypoint sh mikefarah/yq
```

It can be useful to have a bash function to avoid typing the whole docker command:

```bash
yq() {
  docker run --rm -i -v "${PWD}":/workdir mikefarah/yq "$@"
}
```

```bash
yq() {
  podman run --rm -i -v "${PWD}":/workdir mikefarah/yq "$@"
}
```
#### Running as root:

`yq`'s container image no longer runs under root (https://github.com/mikefarah/yq/pull/860). If you'd like to install more things in the container image, or you're having permissions issues when attempting to read/write files you'll need to either:


```
docker run --user="root" -it --entrypoint sh mikefarah/yq
```

```
podman run --user="root" -it --entrypoint sh mikefarah/yq
```

Or, in your Dockerfile:

```
FROM mikefarah/yq

USER root
RUN apk add --no-cache bash
USER yq
```

#### Missing timezone data
By default, the alpine image yq uses does not include timezone data. If you'd like to use the `tz` operator, you'll need to include this data:

```
FROM mikefarah/yq

USER root
RUN apk add --no-cache tzdata
USER yq
```

### GitHub Action
```
  - name: Set foobar to cool
    uses: mikefarah/yq@master
    with:
      cmd: yq -i '.foo.bar = "cool"' 'config.yml'
```

See https://mikefarah.gitbook.io/yq/usage/github-action for more.

### Go Install:
```
go install github.com/mikefarah/yq/v4@latest
```

## Community Supported Installation methods
As these are supported by the community :heart: - however, they may be out of date with the officially supported releases.

# Webi

```
webi yq
```

See [webi](https://webinstall.dev/)
Supported by @adithyasunil26 (https://github.com/webinstall/webi-installers/tree/master/yq)

### Arch Linux

```
pacman -S go-yq
```

### Windows:
[![Chocolatey](https://img.shields.io/chocolatey/v/yq.svg)](https://chocolatey.org/packages/yq)
[![Chocolatey](https://img.shields.io/chocolatey/dt/yq.svg)](https://chocolatey.org/packages/yq)
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
- [Detailed documentation with many examples](https://mikefarah.gitbook.io/yq/)
- Written in portable go, so you can download a lovely dependency free binary
- Uses similar syntax as `jq` but works with YAML, [JSON](https://mikefarah.gitbook.io/yq/usage/convert) and [XML](https://mikefarah.gitbook.io/yq/usage/xml) files
- Fully supports multi document yaml files
- Supports yaml [front matter](https://mikefarah.gitbook.io/yq/usage/front-matter) blocks (e.g. jekyll/assemble)
- Colorized yaml output
- [Date/Time manipulation and formatting with TZ](https://mikefarah.gitbook.io/yq/operators/datetime)
- [Deeply data structures](https://mikefarah.gitbook.io/yq/operators/traverse-read)
- [Sort keys](https://mikefarah.gitbook.io/yq/operators/sort-keys)
- Manipulate yaml [comments](https://mikefarah.gitbook.io/yq/operators/comment-operators), [styling](https://mikefarah.gitbook.io/yq/operators/style), [tags](https://mikefarah.gitbook.io/yq/operators/tag) and [anchors and aliases](https://mikefarah.gitbook.io/yq/operators/anchor-and-alias-operators).
- [Update inplace](https://mikefarah.gitbook.io/yq/v/v4.x/commands/evaluate#flags)
- [Complex expressions to select and update](https://mikefarah.gitbook.io/yq/operators/select#select-and-update-matching-values-in-map)
- Keeps yaml formatting and comments when updating (though there are issues with whitespace)
- [Decode/Encode base64 data](https://mikefarah.gitbook.io/yq/operators/encode-decode)
- [Load content from other files](https://mikefarah.gitbook.io/yq/operators/load)
- [Convert to/from json](https://mikefarah.gitbook.io/yq/v/v4.x/usage/convert)
- [Convert to/from xml](https://mikefarah.gitbook.io/yq/v/v4.x/usage/xml)
- [Convert to/from properties](https://mikefarah.gitbook.io/yq/v/v4.x/usage/properties)
- [Convert to csv/tsv](https://mikefarah.gitbook.io/yq/usage/csv-tsv)
- [General shell completion scripts (bash/zsh/fish/powershell)](https://mikefarah.gitbook.io/yq/v/v4.x/commands/shell-completion)
- [Reduce](https://mikefarah.gitbook.io/yq/operators/reduce) to merge multiple files or sum an array or other fancy things.
- [Github Action](https://mikefarah.gitbook.io/yq/usage/github-action) to use in your automated pipeline (thanks @devorbitus)

## [Usage](https://mikefarah.gitbook.io/yq/)

Check out the [documentation](https://mikefarah.gitbook.io/yq/) for more detailed and advanced usage.

```
Usage:
  yq [flags]
  yq [command]

Examples:

# yq defaults to 'eval' command if no command is specified. See "yq eval --help" for more examples.
yq '.stuff' < myfile.yml # outputs the data at the "stuff" node from "myfile.yml"

yq -i '.stuff = "foo"' myfile.yml # update myfile.yml inplace


Available Commands:
  completion       Generate the autocompletion script for the specified shell
  eval             (default) Apply the expression to each document in each yaml file in sequence
  eval-all         Loads _all_ yaml documents of _all_ yaml files and runs expression once
  help             Help about any command
  shell-completion Generate completion script

Flags:
  -C, --colors                        force print with colors
  -e, --exit-status                   set exit status if there are no matches or null or false is returned
  -f, --front-matter string           (extract|process) first input as yaml front-matter. Extract will pull out the yaml content, process will run the expression against the yaml content, leaving the remaining data intact
      --header-preprocess             Slurp any header comments and separators before processing expression. (default true)
  -h, --help                          help for yq
  -I, --indent int                    sets indent level for output (default 2)
  -i, --inplace                       update the file inplace of first file given.
  -p, --input-format string           [yaml|y|xml|x] parse format for input. Note that json is a subset of yaml. (default "yaml")
  -M, --no-colors                     force print with no colors
  -N, --no-doc                        Don't print document separators (---)
  -n, --null-input                    Don't read input, simply evaluate the expression given. Useful for creating docs from scratch.
  -o, --output-format string          [yaml|y|json|j|props|p|xml|x] output format type. (default "yaml")
  -P, --prettyPrint                   pretty print, shorthand for '... style = ""'
  -s, --split-exp string              print each result (or doc) into a file named (exp). [exp] argument must return a string. You can use $index in the expression as the result counter.
      --unwrapScalar                  unwrap scalar, print the value with no quotes, colors or comments (default true)
  -v, --verbose                       verbose mode
  -V, --version                       Print version information and quit
      --xml-attribute-prefix string   prefix for xml attributes (default "+")
      --xml-content-name string       name for xml content (if no attribute name is present). (default "+content")

Use "yq [command] --help" for more information about a command.
```
## Known Issues / Missing Features
- `yq` attempts to preserve comment positions and whitespace as much as possible, but it does not handle all scenarios (see https://github.com/go-yaml/yaml/tree/v3 for details)
- Powershell has its own...[opinions on quoting yq](https://mikefarah.gitbook.io/yq/usage/tips-and-tricks#quotes-in-windows-powershell)

See [tips and tricks](https://mikefarah.gitbook.io/yq/usage/tips-and-tricks) for more common problems and solutions.
