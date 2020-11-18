The style operator can be used to get or set the style of nodes (e.g. string style, yaml style)
## Examples
### Set tagged style
Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
```
then
```bash
yq eval '.. style="tagged"' sample.yml
```
will output
```yaml
!!map
a: !!str cat
b: !!int 5
c: !!float 3.2
e: !!bool true
```

### Set double quote style
Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
```
then
```bash
yq eval '.. style="double"' sample.yml
```
will output
```yaml
a: "cat"
b: "5"
c: "3.2"
e: "true"
```

### Set single quote style
Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
```
then
```bash
yq eval '.. style="single"' sample.yml
```
will output
```yaml
a: 'cat'
b: '5'
c: '3.2'
e: 'true'
```

### Set literal quote style
Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
```
then
```bash
yq eval '.. style="literal"' sample.yml
```
will output
```yaml
a: |-
  cat
b: |-
  5
c: |-
  3.2
e: |-
  true
```

### Set folded quote style
Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
```
then
```bash
yq eval '.. style="folded"' sample.yml
```
will output
```yaml
a: >-
  cat
b: >-
  5
c: >-
  3.2
e: >-
  true
```

### Set flow quote style
Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
```
then
```bash
yq eval '.. style="flow"' sample.yml
```
will output
```yaml
{a: cat, b: 5, c: 3.2, e: true}
```

### Set empty (default) quote style
Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
```
then
```bash
yq eval '.. style=""' sample.yml
```
will output
```yaml
a: cat
b: 5
c: 3.2
e: true
```

### Set style using a path
Given a sample.yml file of:
```yaml
a: cat
b: double
```
then
```bash
yq eval '.a style=.b' sample.yml
```
will output
```yaml
a: "cat"
b: double
```

### Example 8
Given a sample.yml file of:
```yaml
a: cat
b: dog
```
then
```bash
yq eval '.. style=""' sample.yml
```
will output
```yaml
a: cat
b: dog
```

### Example 9
Given a sample.yml file of:
```yaml
a: cat
b: thing
```
then
```bash
yq eval '.. | style' sample.yml
```
will output
```yaml



```

### Example 10
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq eval '.. | style' sample.yml
```
will output
```yaml


```

