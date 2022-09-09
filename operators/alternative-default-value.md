# Alternative (Default value)

This operator is used to provide alternative (or default) values when a particular expression is either null or false.

{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

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

## Update or create - entity exists
This initialises `a` if it's not present

Given a sample.yml file of:
```yaml
a: 1
```
then
```bash
yq '(.a // (.a = 0)) += 1' sample.yml
```
will output
```yaml
a: 2
```

## Update or create - entity does not exist
This initialises `a` if it's not present

Given a sample.yml file of:
```yaml
b: camel
```
then
```bash
yq '(.a // (.a = 0)) += 1' sample.yml
```
will output
```yaml
b: camel
a: 1
```

