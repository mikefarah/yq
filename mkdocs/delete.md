```
yq d <yaml_file> <path_to_delete>
```

### To Stdout
Given a sample.yaml file of:
```yaml
b:
  c: 2
  apples: green
```
then
```bash
yq d sample.yaml b.c
```
will output:
```yaml
b:
  apples: green
```

### From STDIN
```bash
cat sample.yaml | yq d - b.c
```

### Deleting array elements
Given a sample.yaml file of:
```yaml
b:
  c: 
    - 1
    - 2
    - 3
```
then
```bash
yq d sample.yaml 'b.c[1]'
```
will output:
```yaml
b:
  c:
  - 1
  - 3
```

### Deleting nodes in-place
Given a sample.yaml file of:
```yaml
b:
  c: 2
  apples: green
```
then
```bash
yq d -i sample.yaml b.c
```
will update the sample.yaml file so that the 'c' node is deleted



### Multiple Documents - delete from single document
Given a sample.yaml file of:
```yaml
something: else
field: leaveMe
---
b:
  c: 2
field: deleteMe
```
then
```bash
yq w -d1 sample.yaml field
```
will output:
```yaml
something: else
field: leaveMe
---
b:
  c: 2
```

### Multiple Documents - delete from all documents
Given a sample.yaml file of:
```yaml
something: else
field: deleteMe
---
b:
  c: 2
field: deleteMeToo
```
then
```bash
yq w -d'*' sample.yaml field
```
will output:
```yaml
something: else
---
b:
  c: 2
```

Note that '*' is in quotes to avoid being interpreted by your shell.

{!snippets/keys_with_dots.md!}
