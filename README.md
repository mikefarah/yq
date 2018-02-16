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
or, [Download latest binary](https://github.com/mikefarah/yq/releases/latest) or alternatively:
```
go get github.com/mikefarah/yq
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

## [Usage](http://mikefarah.github.io/yq/)

Check out the [documentation](http://mikefarah.github.io/yq/) for more detailed and advanced usage.

```
Usage:
  yq [flags]
  yq [command]

Available Commands:
  help        Help about any command
  merge       yq m [--inplace/-i] [--overwrite/-x] sample.yaml sample2.yaml
  new         yq n [--script/-s script_file] a.b.c newValueForC
  read        yq r sample.yaml a.b.c
  write       yq w [--inplace/-i] [--script/-s script_file] sample.yaml a.b.c newValueForC

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
3. apply changes
4. `make [local] build`
5. profit
