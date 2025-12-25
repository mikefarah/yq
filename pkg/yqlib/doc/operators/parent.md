# Parent

Parent simply returns the parent nodes of the matching nodes.

## Simple example
Given a sample.yml file of:
```yaml
a:
  nested: cat
```
then
```bash
yq '.a.nested | parent' sample.yml
```
will output
```yaml
nested: cat
```

## Parent of nested matches
Given a sample.yml file of:
```yaml
a:
  fruit: apple
  name: bob
b:
  fruit: banana
  name: sam
```
then
```bash
yq '.. | select(. == "banana") | parent' sample.yml
```
will output
```yaml
fruit: banana
name: sam
```

## Get parent attribute
Given a sample.yml file of:
```yaml
a:
  fruit: apple
  name: bob
b:
  fruit: banana
  name: sam
```
then
```bash
yq '.. | select(. == "banana") | parent.name' sample.yml
```
will output
```yaml
sam
```

## Get parents
Match all parents

Given a sample.yml file of:
```yaml
a:
  b:
    c: cat
```
then
```bash
yq '.a.b.c | parents' sample.yml
```
will output
```yaml
- c: cat
- b:
    c: cat
- a:
    b:
      c: cat
```

## Get the top (root) parent
Use negative numbers to get the top parents. You can think of this as indexing into the 'parents' array above

Given a sample.yml file of:
```yaml
a:
  b:
    c: cat
```
then
```bash
yq '.a.b.c | parent(-1)' sample.yml
```
will output
```yaml
a:
  b:
    c: cat
```

## Root
Alias for parent(-1), returns the top level parent. This is usually the document node.

Given a sample.yml file of:
```yaml
a:
  b:
    c: cat
```
then
```bash
yq '.a.b.c | root' sample.yml
```
will output
```yaml
a:
  b:
    c: cat
```

## N-th parent
You can optionally supply the number of levels to go up for the parent, the default being 1.

Given a sample.yml file of:
```yaml
a:
  b:
    c: cat
```
then
```bash
yq '.a.b.c | parent(2)' sample.yml
```
will output
```yaml
b:
  c: cat
```

## N-th parent - another level
Given a sample.yml file of:
```yaml
a:
  b:
    c: cat
```
then
```bash
yq '.a.b.c | parent(3)' sample.yml
```
will output
```yaml
a:
  b:
    c: cat
```

## N-th negative
Similarly, use negative numbers to index backwards from the parents array

Given a sample.yml file of:
```yaml
a:
  b:
    c: cat
```
then
```bash
yq '.a.b.c | parent(-2)' sample.yml
```
will output
```yaml
b:
  c: cat
```

## No parent
Given a sample.yml file of:
```yaml
{}
```
then
```bash
yq 'parent' sample.yml
```
will output
```yaml
```

