Explodes (or dereferences) aliases and anchors.
## Examples
### Explode alias and anchor
Given a sample.yml file of:
```yaml
f:
  a: &a cat
  b: *a
```
then
```bash
yq eval 'explode(.f)' sample.yml
```
will output
```yaml
{f: {a: cat, b: cat}}
```

### Explode with no aliases or anchors
Given a sample.yml file of:
```yaml
a: mike
```
then
```bash
yq eval 'explode(.a)' sample.yml
```
will output
```yaml
a: mike
```

### Explode with alias keys
Given a sample.yml file of:
```yaml
f:
  a: &a cat
  *a: b
```
then
```bash
yq eval 'explode(.f)' sample.yml
```
will output
```yaml
{f: {a: cat, cat: b}}
```

### Explode with merge anchors
Given a sample.yml file of:
```yaml
foo: &foo
  a: foo_a
  thing: foo_thing
  c: foo_c
bar: &bar
  b: bar_b
  thing: bar_thing
  c: bar_c
foobarList:
  b: foobarList_b
  !!merge <<:
    - *foo
    - *bar
  c: foobarList_c
foobar:
  c: foobar_c
  !!merge <<: *foo
  thing: foobar_thing
```
then
```bash
yq eval 'explode(.)' sample.yml
```
will output
```yaml
foo:
  a: foo_a
  thing: foo_thing
  c: foo_c
bar:
  b: bar_b
  thing: bar_thing
  c: bar_c
foobarList:
  b: bar_b
  a: foo_a
  thing: bar_thing
  c: foobarList_c
foobar:
  c: foo_c
  a: foo_a
  thing: foobar_thing
```

