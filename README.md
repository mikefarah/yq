# yaml [![Build Status](https://travis-ci.org/mikefarah/yaml.svg?branch=master)](https://travis-ci.org/mikefarah/yaml)
yaml is a lightweight and portable command-line YAML processor

The aim of the project is to be the [jq](https://github.com/stedolan/jq) or sed of yaml files.

## Install
[Download latest binary](https://github.com/mikefarah/yaml/releases/latest) or alternatively:
```
go get github.com/mikefarah/yaml
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

## [Usage](http://mikefarah.github.io/yaml/)

Check out the [documentation](http://mikefarah.github.io/yaml/) for more detailed and advanced usage.

```
Usage:
  yaml [flags]
  yaml [command]

Available Commands:
  help        Help about any command
  merge       yaml m [--inplace/-i] [--overwrite/-x] sample.yaml sample2.yaml
  new         yaml n [--script/-s script_file] a.b.c newValueForC
  read        yaml r sample.yaml a.b.c
  write       yaml w [--inplace/-i] [--script/-s script_file] sample.yaml a.b.c newValueForC

Flags:
  -h, --help      help for yaml
  -j, --tojson    output as json
  -t, --trim      trim yaml output (default true)
  -v, --verbose   verbose mode
  -V, --version   Print version information and quit

Use "yaml [command] --help" for more information about a command.
```

## Contribute
1. `make [local] vendor` OR run `govendor sync` [link](https://github.com/kardianos/govendor)
2. add unit tests
3. apply changes
4. `make`
5. profit
