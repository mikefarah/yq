# Column

Returns the column of the matching node. Starts from 1, 0 indicates there was no column data.

## Returns column of _value_ node
Given a sample.yml file of:
```yaml
a: cat
b: bob
```
then
```bash
yq '.b | column' sample.yml
```
will output
```yaml
4
```

## Returns column of _key_ node
Pipe through the key operator to get the column of the key

Given a sample.yml file of:
```yaml
a: cat
b: bob
```
then
```bash
yq '.b | key | column' sample.yml
```
will output
```yaml
1
```

## First column is 1
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq '.a | key | column' sample.yml
```
will output
```yaml
1
```

## No column data is 0
Running
```bash
yq --null-input '{"a": "new entry"} | column'
```
will output
```yaml
0
```

