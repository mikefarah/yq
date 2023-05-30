# Flatten
This recursively flattens arrays.

## Flatten
Recursively flattens all arrays

Given a sample.yml file of:
```yaml
- 1
- - 2
- - - 3
```
then
```bash
yq 'flatten' sample.yml
```
will output
```yaml
- 1
- 2
- 3
```

## Flatten with depth of one
Given a sample.yml file of:
```yaml
- 1
- - 2
- - - 3
```
then
```bash
yq 'flatten(1)' sample.yml
```
will output
```yaml
- 1
- 2
- - 3
```

## Flatten empty array
Given a sample.yml file of:
```yaml
- []
```
then
```bash
yq 'flatten' sample.yml
```
will output
```yaml
[]
```

## Flatten array of objects
Given a sample.yml file of:
```yaml
- foo: bar
- - foo: baz
```
then
```bash
yq 'flatten' sample.yml
```
will output
```yaml
- foo: bar
- foo: baz
```

