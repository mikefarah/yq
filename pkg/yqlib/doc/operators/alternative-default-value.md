# Alternative (Default value)

This operator is used to provide alternative (or default) values when a particular expression is either null or false.

## LHS is defined
Given a sample.yml file of:
```yaml
a: bridge
```
then
```bash
yq '.a // "hello"' sample.yml
```
will output
```yaml
bridge
```

## LHS is not defined
Given a sample.yml file of:
```yaml
{}
```
then
```bash
yq '.a // "hello"' sample.yml
```
will output
```yaml
hello
```

## LHS is null
Given a sample.yml file of:
```yaml
a: ~
```
then
```bash
yq '.a // "hello"' sample.yml
```
will output
```yaml
hello
```

## LHS is false
Given a sample.yml file of:
```yaml
a: false
```
then
```bash
yq '.a // "hello"' sample.yml
```
will output
```yaml
hello
```

## RHS is an expression
Given a sample.yml file of:
```yaml
a: false
b: cat
```
then
```bash
yq '.a // .b' sample.yml
```
will output
```yaml
cat
```

