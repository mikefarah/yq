```
yq delete <yaml_file|-> <path_expression>
```

The delete command will delete all the matching nodes for the path expression in the given yaml input.

See docs for [path expression](path_expressions.md) for more details.


## Deleting from a simple document
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
will output
```yaml
b:
  apples: green
```

## From STDIN
Use "-" (without quotes) in-place of a file name if you wish to pipe in input from STDIN.

```bash
cat sample.yaml | yq d - b.c
```

## Deleting in-place
```bash
yq d -i sample.yaml b.c
```
will update the sample.yaml file so that the 'c' node is deleted


## Multiple Documents

### Delete from single document
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

### Delete from all documents
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
