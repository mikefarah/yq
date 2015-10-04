# yaml
yaml is portable command line tool written in go

Allows you to read (and soon update) yaml files given a yaml path. All in a lovely dependency free binary!

[Download latest release](https://github.com/mikefarah/yaml/releases/latest)

or alternatively install using go get:
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

### Handling '.' in the yaml key
Given a sample.yaml file of:
```yaml
b.x:
  c: 2
```
then
```bash
yaml sample.yaml \"b.x\".c
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
yaml sample.yaml b.e[1].name
```
will output 'sam'

### Updating yaml
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
