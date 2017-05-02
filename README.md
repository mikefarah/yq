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

## [Usage](http://mikefarah.github.io/yaml/)

Check out the [documentation](http://mikefarah.github.io/yaml/) for more detailed and advanced usage.

```
Usage:
  yaml [command]

Available Commands:
  read        yaml r sample.yaml a.b.c
  write       yaml w [--inplace/-i] [--script/-s script_file] sample.yaml a.b.c newValueForC
  new         yaml n [--script/-s script_file] a.b.c newValueForC

Flags:
  -h, --help[=false]: help for yaml
  -j, --tojson[=false]: output as json
  -t, --trim[=true]: trim yaml output
  -v, --verbose[=false]: verbose mode

Use "yaml [command] --help" for more information about a command.
```

## Contribute
1. run `govendor sync` [link](https://github.com/kardianos/govendor)
2. add unit tests
3. make changes
4. run ./precheckin.sh
5. profit
