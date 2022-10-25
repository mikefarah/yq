# Collect into Array

This creates an array using the expression between the square brackets.


## Collect empty
Running
```bash
yq --null-input '[]'
```
will output
```yaml
[]
```

## Collect single
Running
```bash
yq --null-input '["cat"]'
```
will output
```yaml
- cat
```

## Collect many
Given a sample.yml file of:
```yaml
a: cat
b: dog
```
then
```bash
yq '[.a, .b]' sample.yml
```
will output
```yaml
- cat
- dog
```

