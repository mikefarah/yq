# Pipe

Pipe the results of an expression into another. Like the bash operator.

## Simple Pipe
Given a sample.yml file of:
```yaml
a:
  b: cat
```
then
```bash
yq eval '.a | .b' sample.yml
```
will output
```yaml
cat
```

## Multiple updates
Given a sample.yml file of:
```yaml
a: cow
b: sheep
c: same
```
then
```bash
yq eval '.a = "cat" | .b = "dog"' sample.yml
```
will output
```yaml
a: cat
b: dog
c: same
```

