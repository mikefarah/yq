
## Get kind
Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
f: []
g: {}
h: null
```
then
```bash
yq '.. | kind' sample.yml
```
will output
```yaml
map
scalar
scalar
scalar
scalar
seq
map
scalar
```

## Get kind, ignores custom tags
Unlike tag, kind is not affected by custom tags.

Given a sample.yml file of:
```yaml
a: !!thing cat
b: !!foo {}
c: !!bar []
```
then
```bash
yq '.. | kind' sample.yml
```
will output
```yaml
map
scalar
map
seq
```

## Add comments only to scalars
An example of how you can use kind

Given a sample.yml file of:
```yaml
a:
  b: 5
  c: 3.2
e: true
f: []
g: {}
h: null
```
then
```bash
yq '(.. | select(kind == "scalar")) line_comment = "this is a scalar"' sample.yml
```
will output
```yaml
a:
  b: 5 # this is a scalar
  c: 3.2 # this is a scalar
e: true # this is a scalar
f: []
g: {}
h: null # this is a scalar
```

