Updates the LHS using the expression on the RHS. Note that the RHS runs against the _original_ LHS value, so that you can evaluate a new value based on the old (e.g. increment).
## Examples
### Update parent to be the child value
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
{a: {g: foof}}
```

### Update string value
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
{a: {b: frog}}
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
{a: {b: frog, c: cactus}}
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
[bogs, apple, bogs]
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
{a: {b: bogs}}
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
{a: {b: [bogs]}}
```

