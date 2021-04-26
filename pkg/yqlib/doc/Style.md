The style operator can be used to get or set the style of nodes (e.g. string style, yaml style)
## Update and set style of a particular node (simple)
Given a sample.yml file of:
```yaml
a:
  b: thing
  c: something
```
then
```bash
yq eval '.a.b = "new" | .a.b style="double"' sample.yml
```
will output
```yaml
a:
  b: "new"
  c: something
```

## Update and set style of a particular node using path variables
You can use a variable to re-use a path

Given a sample.yml file of:
```yaml
a:
  b: thing
  c: something
```
then
```bash
yq eval '.a.b as $x | $x = "new" | $x style="double"' sample.yml
```
will output
```yaml
a:
  b: "new"
  c: something
```

## Set tagged style
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

## Set double quote style
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

## Set double quote style on map keys too
Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
```
then
```bash
yq eval '... style="double"' sample.yml
```
will output
```yaml
"a": "cat"
"b": "5"
"c": "3.2"
"e": "true"
```

## Set single quote style
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

## Set literal quote style
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

## Set folded quote style
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

## Set flow quote style
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

## Reset style - or pretty print
Set empty (default) quote style, note the usage of `...` to match keys too. Note that there is a `--prettyPrint/-P` short flag for this.

Given a sample.yml file of:
```yaml
a: cat
"b": 5
'c': 3.2
"e": true
```
then
```bash
yq eval '... style=""' sample.yml
```
will output
```yaml
a: cat
b: 5
c: 3.2
e: true
```

## Set style relatively with assign-update
Given a sample.yml file of:
```yaml
a: single
b: double
```
then
```bash
yq eval '.[] style |= .' sample.yml
```
will output
```yaml
a: 'single'
b: "double"
```

## Read style
Given a sample.yml file of:
```yaml
{a: "cat", b: 'thing'}
```
then
```bash
yq eval '.. | style' sample.yml
```
will output
```yaml
flow
double
single
```

