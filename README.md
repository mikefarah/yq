# yq [![Build Status](https://travis-ci.org/mikefarah/yq.svg?branch=master)](https://travis-ci.org/mikefarah/yq)
a lightweight and portable command-line YAML processor

The aim of the project is to be the [jq](https://github.com/stedolan/jq) or sed of yaml files.

## Install
On MacOS:
```
brew install yq
```
On Ubuntu and other Linux distros supporting `snap` packages:
```
snap install yq
```
On Ubuntu 16.04 or higher from Debian package:
```
sudo add-apt-repository ppa:rmescandon/yq
sudo apt update
sudo apt install yq -y
```
or, [Download latest binary](https://github.com/mikefarah/yq/releases/latest) or alternatively:
```
go get gopkg.in/mikefarah/yq.v2
```

## Run with Docker

Oneshot use:

```bash
docker run -v ${PWD}:/workdir mikefarah/yq yq [flags] <command> FILE...
```

Run commands interactively:

```bash
docker run -it -v ${PWD}:/workdir mikefarah/yq sh
```

## Features
- Written in portable go, so you can download a lovely dependency free binary
- Deep read a yaml file with a given path
- Update a yaml file given a path
- Update a yaml file given a script file
- Update creates any missing entries in the path on the fly
- Create a yaml file given a deep path and value
- Create a yaml file given a script file
- Convert from json to yaml
- Convert from yaml to json
- Pipe data in by using '-'
- Merge multiple yaml files where each additional file sets values for missing or null value keys.
- Merge multiple yaml files with overwrite to support overriding previous values.
- Supports multiple documents in a single yaml file

## [Usage](http://mikefarah.github.io/yq/)

Check out the [documentation](http://mikefarah.github.io/yq/) for more detailed and advanced usage.

```
Usage:
  yq [flags]
  yq [command]

Available Commands:
  delete      yq d [--inplace/-i] [--doc/-d document_index] sample.yaml a.b.c
  help        Help about any command
  merge       yq m [--inplace/-i] [--doc/-d document_index] [--overwrite/-x] sample.yaml sample2.yaml
  new         yq n [--script/-s script_file] a.b.c newValueForC
  read        yq r [--doc/-d document_index] sample.yaml a.b.c
  write       yq w [--inplace/-i] [--script/-s script_file] [--doc/-d document_index] sample.yaml a.b.c newValueForC

Flags:
  -h, --help      help for yq
  -j, --tojson    output as json
  -t, --trim      trim yaml output (default true)
  -v, --verbose   verbose mode
  -V, --version   Print version information and quit

Use "yq [command] --help" for more information about a command.
```

## Contribute
1. `make [local] vendor`
2. add unit tests
3. apply changes (use govendor with a preference to [gopkg](https://gopkg.in/) for package dependencies)
4. `make [local] build`
5. If required, update the user documentation 
    - Update README.md and/or documentation under the mkdocs folder
    - `make [local] build-docs`
    - browse to docs/index.html and check your changes 
6. profit
