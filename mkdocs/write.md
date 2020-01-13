```
yq w <yaml_file> <path_expression> <new value>
```

Updates all the matching nodes of path expression to the supplied value.

See docs for [path expression](path_expressions.md) for more details.

## Basic
Given a sample.yaml file of:
```yaml
b:
  c: 2
```
then
```bash
yq w sample.yaml b.c cat
```
will output:
```yaml
b:
  c: cat
```

### Updating files in-place
```bash
yq w -i sample.yaml b.c cat
```
will update the sample.yaml file so that the value of 'c' is cat.

## From STDIN
```bash
cat sample.yaml | yq w - b.c blah
```

## Adding new fields
Any missing fields in the path will be created on the fly.

Given a sample.yaml file of:
```yaml
b:
  c: 2
```
then
```bash
yq w sample.yaml b.d[+] "new thing"
```
will output:
```yaml
b:
  c: cat
  d:
    - new thing
```

## Appending value to an array field
Given a sample.yaml file of:
```yaml
b:
  c: 2
  d:
    - new thing
    - foo thing
```
then
```bash
yq w sample.yaml "b.d[+]" "bar thing"
```
will output:
```yaml
b:
  c: cat
  d:
    - new thing
    - foo thing
    - bar thing
```

Note that the path is in quotes to avoid the square brackets being interpreted by your shell.

## Multiple Documents 
### Update a single document
Given a sample.yaml file of:
```yaml
something: else
---
b:
  c: 2
```
then
```bash
yq w -d1 sample.yaml b.c 5
```
will output:
```yaml
something: else
---
b:
  c: 5
```

### Update all documents
Given a sample.yaml file of:
```yaml
something: else
---
b:
  c: 2
```
then
```bash
yq w -d'*' sample.yaml b.c 5
```
will output:
```yaml
something: else
b:
  c: 5
---
b:
  c: 5
```

UPDATE THIS
UPDATE THIS 
INCLUDE DELETE EXAMPLE

## Updating multiple values with a script 
Given a sample.yaml file of:
```yaml
b:
  c: 2
  e:
    - name: Billy Bob
```
and a script update_instructions.yaml of:
```yaml
b.c: 3
b.e[+].name: Howdy Partner
```
then

```bash
yq w -s update_instructions.yaml sample.yaml
```
will output:
```yaml
b:
  c: 3
  e:
    - name: Howdy Partner
```

And, of course, you can pipe the instructions in using '-':
```bash
cat update_instructions.yaml | yq w -s - sample.yaml
```
