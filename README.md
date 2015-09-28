# yaml
yaml command line tool written in go

Allows you to read (and soon update) yaml files given a yaml path.

## Install
```
go get github.com/mikefarah/yaml
```

## Read examples
```
yaml <yaml file> <path>
```

### Basic
Given a sample.yaml file of:
```yaml
b:
  c: 2
```
then
```bash
yaml sample.yaml b.c
```
will output the value of '2'.

### Arrays
You can give an index to access a specific element:
e.g.: given a sample file of
```yaml
b:
  e:
    - name: fred
      value: 3
    - name: sam
      value: 4
```
then
```
yaml sample.yaml b.e.1.name
```
will output 'sam'

## TODO
* Updating yaml files
* Handling '.' in path names
