This operator is used to update node values. It can be used in either the:

### plain form: `=`
Which will assign the LHS node values to the RHS node values. The RHS expression is run against the matching nodes in the pipeline.

### relative form: `|=`
This will do a similar thing to the plain form, however, the RHS expression is run against _the LHS nodes_. This is useful for updating values based on old values, e.g. increment.
## Create yaml file
Running
```bash
yq eval --null-input '.a.b = "cat" | .x = "frog"'
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
a: {b: {g: foof}}
'': null
```
then
```bash
yq eval '.a |= .b' sample.yml
```
will output
```yaml
a: {g: foof}
'': null
```

## Update node to be the sibling value
Given a sample.yml file of:
```yaml
a: {b: child}
b: sibling
'': null
```
then
```bash
yq eval '.a = .b' sample.yml
```
will output
```yaml
a: sibling
b: sibling
'': null
```

## Updated multiple paths
Given a sample.yml file of:
```yaml
a: fieldA
b: fieldB
c: fieldC
'': null
```
then
```bash
yq eval '(.a, .c) |= "potatoe"' sample.yml
```
will output
```yaml
a: potatoe
b: fieldB
c: potatoe
'': null
```

## Update string value
Given a sample.yml file of:
```yaml
a: {b: apple}
'': null
```
then
```bash
yq eval '.a.b = "frog"' sample.yml
```
will output
```yaml
a: {b: frog}
'': null
```

## Update string value via |=
Note there is no difference between `=` and `|=` when the RHS is a scalar

Given a sample.yml file of:
```yaml
a: {b: apple}
'': null
```
then
```bash
yq eval '.a.b |= "frog"' sample.yml
```
will output
```yaml
a: {b: frog}
'': null
```

## Update selected results
Given a sample.yml file of:
```yaml
a: {b: apple, c: cactus}
'': null
```
then
```bash
yq eval '(.a[] | select(. == "apple")) = "frog"' sample.yml
```
will output
```yaml
a: {b: frog, c: cactus}
'': null
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
yq eval '(.[] | select(. == "*andy")) = "bogs"' sample.yml
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
yq eval '.a.b |= "bogs"' sample.yml
```
will output
```yaml
{a: {b: bogs}}
```

## Update empty object and array
Given a sample.yml file of:
```yaml
{}
```
then
```bash
yq eval '.a.b.[0] |= "bogs"' sample.yml
```
will output
```yaml
{a: {b: [bogs]}}
```

