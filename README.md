# yaml
yaml is portable command line tool written in go

Allows you to read and update yaml files from bash (or whatever). All in a lovely dependency free binary!

[Download latest release](https://github.com/mikefarah/yaml/releases/latest)

or alternatively install using go get:
```
go get github.com/mikefarah/yaml
```

## Read examples
```
yaml r <yaml file> <path>
```

### Basic
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

### Reading from STDIN
Given a sample.yaml file of:
```bash
cat sample.yaml | yaml r - b.c
```
will output the value of '2'.

### Splat
Given a sample.yaml file of:
```yaml
---
bob:
  item1:
    cats: bananas
  item2:
    cats: apples
```
then
```bash
yaml r sample.yaml bob.*.cats
```
will output
```yaml
- bananas
- apples
```

### Handling '.' in the yaml key
Given a sample.yaml file of:
```yaml
b.x:
  c: 2
```
then
```bash
yaml r sample.yaml \"b.x\".c
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
yaml r sample.yaml b.e[1].name
```
will output 'sam'

### Array Splat
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
yaml r sample.yaml b.e[*].name
```
will output:
```
- fred
- sam
```

## Update examples

### Update to stdout
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

### Updating yaml in-place
Given a sample.yaml file of:
```yaml
b:
  c: 2
```
then
```bash
yaml w -i sample.yaml b.c cat
```
will update the sample.yaml file so that the value of 'c' is cat.
