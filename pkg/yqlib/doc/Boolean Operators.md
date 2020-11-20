The `or` and `and` operators take two parameters and return a boolean result. These are most commonly used with the `select` operator to filter particular nodes.
## Examples
### Update node to be the child value
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

### Update node to be the sibling value
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

### Updated multiple paths
Given a sample.yml file of:
```yaml
a: fieldA
b: fieldB
c: fieldC
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
```

### Update string value
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

### Update string value via |=
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

### Update selected results
Given a sample.yml file of:
```yaml
a:
  b: apple
  c: cactus
```
then
```bash
yq eval '.a[] | select(. == "apple") |= "frog"' sample.yml
```
will output
```yaml
a:
  b: frog
  c: cactus
```

### Update array values
Given a sample.yml file of:
```yaml
- candy
- apple
- sandy
```
then
```bash
yq eval '.[] | select(. == "*andy") |= "bogs"' sample.yml
```
will output
```yaml
- bogs
- apple
- bogs
```

### Update empty object
Given a sample.yml file of:
```yaml
'': null
```
then
```bash
yq eval '.a.b |= "bogs"' sample.yml
```
will output
```yaml
'': null
a:
  b: bogs
```

### Update empty object and array
Given a sample.yml file of:
```yaml
'': null
```
then
```bash
yq eval '.a.b[0] |= "bogs"' sample.yml
```
will output
```yaml
'': null
a:
  b:
    - bogs
```

