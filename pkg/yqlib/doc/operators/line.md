# Line

Returns the line of the matching node. Starts from 1, 0 indicates there was no line data.

## Returns line of _value_ node
Given a sample.yml file of:
```yaml
a: cat
b:
  c: cat
```
then
```bash
yq '.b | line' sample.yml
```
will output
```yaml
3
```

## Returns line of _key_ node
Pipe through the key operator to get the line of the key

Given a sample.yml file of:
```yaml
a: cat
b:
  c: cat
```
then
```bash
yq '.b | key| line' sample.yml
```
will output
```yaml
2
```

## First line is 1
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq '.a | line' sample.yml
```
will output
```yaml
1
```

## No line data is 0
Running
```bash
yq --null-input '{"a": "new entry"} | line'
```
will output
```yaml
0
```

