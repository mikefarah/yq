# Assign (Update)

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
a:
  b:
    g: foof
```

then

```bash
yq eval '.a |= .b' sample.yml
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
yq eval '.[] |= . * 2' sample.yml
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
yq eval '.a = .b' sample.yml
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
yq eval '(.a, .c) = "potatoe"' sample.yml
```

will output

```yaml
a: potatoe
b: fieldB
c: potatoe
```

## Update string value

Given a sample.yml file of:

```yaml
a:
  b: apple
```

then

```bash
yq eval '.a.b = "frog"' sample.yml
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
yq eval '.a.b |= "frog"' sample.yml
```

will output

```yaml
a:
  b: frog
```

## Update selected results

Given a sample.yml file of:

```yaml
a:
  b: apple
  c: cactus
```

then

```bash
yq eval '(.a[] | select(. == "apple")) = "frog"' sample.yml
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
