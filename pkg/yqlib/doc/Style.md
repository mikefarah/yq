The style operator can be used to get or set the style of nodes (e.g. string style, yaml style)
## Set tagged style
Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
'': null
```
then
```bash
yq eval '.. style="tagged"' sample.yml
```
will output
```yaml
!!map
a: cat
b: 5
c: 3.2
e: true
'': !!null null
```

## Set double quote style
Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
'': null
```
then
```bash
yq eval '.. style="double"' sample.yml
```
will output
```yaml
a: cat
b: 5
c: 3.2
e: true
'': "null"
```

## Set single quote style
Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
'': null
```
then
```bash
yq eval '.. style="single"' sample.yml
```
will output
```yaml
a: cat
b: 5
c: 3.2
e: true
'': 'null'
```

## Set literal quote style
Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
'': null
```
then
```bash
yq eval '.. style="literal"' sample.yml
```
will output
```yaml
a: cat
b: 5
c: 3.2
e: true
'': |-
  null
```

## Set folded quote style
Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
'': null
```
then
```bash
yq eval '.. style="folded"' sample.yml
```
will output
```yaml
a: cat
b: 5
c: 3.2
e: true
'': >-
  null
```

## Set flow quote style
Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
'': null
```
then
```bash
yq eval '.. style="flow"' sample.yml
```
will output
```yaml
{a: cat, b: 5, c: 3.2, e: true, '': null}
```

## Pretty print
Set empty (default) quote style

Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
'': null
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
'': null
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

```

