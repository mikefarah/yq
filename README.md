# yaml [![Build Status](https://travis-ci.org/mikefarah/yaml.svg?branch=master)](https://travis-ci.org/mikefarah/yaml)
yaml is a lightweight and flexible command-line YAML processor

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

### Read
```
yaml r <yaml file> <path>
```

Given a sample.yaml file of:
```yaml
b:
  c: 2
```
then
```bash
yaml r sample.yaml b.c
```
will output the value of '2'.


## Update
Existing yaml files can be updated via the write command

Given a sample.yaml file of:
```yaml
b:
  c: 2
```
then
```bash
yaml w sample.yaml b.c cat
```
will output:
```yaml
b:
  c: cat
```

## Create
Yaml files can be created using the 'new' command. This works in the same way as the write command, but you don't pass in an existing Yaml file.

### Creating a simple yaml file
```bash
yaml n b.c cat
```
will output:
```yaml
b:
  c: cat
```

## Converting to and from json

### Yaml2json
To convert output to json, use the --tojson (or -j) flag. This can be used with any command.

### json2yaml
To read in json, use the --fromjson (or -J) flag. This can be used with any command.
