Paths can be prefixed using the 'prefix' command.
The complete yaml content will be nested inside the new prefix path.

```
yq p <yaml_file> <path>
```

### To Stdout
Given a data1.yaml file of:
```yaml
a: simple
b: [1, 2]
```
then
```bash
yq p data1.yaml c
```
will output:
```yaml
c:
  a: simple
  b: [1, 2]
```

### Arbitrary depth
Given a data1.yaml file of:
```yaml
a:
  b: [1, 2]
```
then
```bash
yq p data1.yaml c.d
```
will output:
```yaml
c:
  d:
    a:
      b: [1, 2]
```

### Updating files in-place
Given a data1.yaml file of:
```yaml
a: simple
b: [1, 2]
```
then
```bash
yq p -i data1.yaml c
```
will update the data1.yaml file so that the path 'c' is prefixed to all other paths.

### Multiple Documents - update a single document
Given a data1.yaml file of:
```yaml
something: else
---
a: simple
b: cat
```
then
```bash
yq p -d1 data1.yaml c
```
will output:
```yaml
something: else
---
c:
  a: simple
  b: cat
```

### Multiple Documents - update a single document
Given a data1.yaml file of:
```yaml
something: else
---
a: simple
b: cat
```
then
```bash
yq p -d'*' data1.yaml c
```
will output:
```yaml
c:
  something: else
---
c:
  a: simple
  b: cat
```
