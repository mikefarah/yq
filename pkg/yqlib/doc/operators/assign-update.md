# Assign (Update)

This operator is used to update node values. It can be used in either the:

### plain form: `=`
Which will assign the LHS node values to the RHS node values. The RHS expression is run against the matching nodes in the pipeline.

### relative form: `|=`
This will do a similar thing to the plain form, however, the RHS expression is run against _the LHS nodes_. This is useful for updating values based on old values, e.g. increment.
{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Create yaml file
Running
```bash
yq --null-input '.a.b = "cat" | .x = "frog"'
```
will output
```yaml
a:
  b: cat
x: frog
```

## Update node to be the child value
Given a sample.yml file of:
```yaml
a:
  b:
    g: foof
```
then
```bash
yq '.a |= .b' sample.yml
```
will output
```yaml
a:
  g: foof
```

## Double elements in an array
Given a sample.yml file of:
```yaml
- 1
- 2
- 3
```
then
```bash
yq '.[] |= . * 2' sample.yml
```
will output
```yaml
- 2
- 4
- 6
```

## Update node from another file
Note this will also work when the second file is a scalar (string/number)

Given a sample.yml file of:
```yaml
a: apples
```
And another sample another.yml file of:
```yaml
b: bob
```
then
```bash
yq eval-all 'select(fileIndex==0).a = select(fileIndex==1) | select(fileIndex==0)' sample.yml another.yml
```
will output
```yaml
a:
  b: bob
```

## Update node to be the sibling value
Given a sample.yml file of:
```yaml
a:
  b: child
b: sibling
```
then
```bash
yq '.a = .b' sample.yml
```
will output
```yaml
a: sibling
b: sibling
```

## Updated multiple paths
Given a sample.yml file of:
```yaml
a: fieldA
b: fieldB
c: fieldC
```
then
```bash
yq '(.a, .c) = "potato"' sample.yml
```
will output
```yaml
a: potato
b: fieldB
c: potato
```

## Update string value
Given a sample.yml file of:
```yaml
a:
  b: apple
```
then
```bash
yq '.a.b = "frog"' sample.yml
```
will output
```yaml
a:
  b: frog
```

## Update string value via |=
Note there is no difference between `=` and `|=` when the RHS is a scalar

Given a sample.yml file of:
```yaml
a:
  b: apple
```
then
```bash
yq '.a.b |= "frog"' sample.yml
```
will output
```yaml
a:
  b: frog
```

## Update deeply selected results
Note that the LHS is wrapped in brackets! This is to ensure we don't first filter out the yaml and then update the snippet.

Given a sample.yml file of:
```yaml
a:
  b: apple
  c: cactus
```
then
```bash
yq '(.a[] | select(. == "apple")) = "frog"' sample.yml
```
will output
```yaml
a:
  b: frog
  c: cactus
```

## Update array values
Given a sample.yml file of:
```yaml
- candy
- apple
- sandy
```
then
```bash
yq '(.[] | select(. == "*andy")) = "bogs"' sample.yml
```
will output
```yaml
- bogs
- apple
- bogs
```

## Update empty object
Given a sample.yml file of:
```yaml
{}
```
then
```bash
yq '.a.b |= "bogs"' sample.yml
```
will output
```yaml
a:
  b: bogs
```

## Update node value that has an anchor
Anchor will remaple

Given a sample.yml file of:
```yaml
a: &cool cat
```
then
```bash
yq '.a = "dog"' sample.yml
```
will output
```yaml
a: &cool dog
```

## Update empty object and array
Given a sample.yml file of:
```yaml
{}
```
then
```bash
yq '.a.b.[0] |= "bogs"' sample.yml
```
will output
```yaml
a:
  b:
    - bogs
```

