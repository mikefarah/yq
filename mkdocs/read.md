```
yq r <yaml_file|json_file> <path_expression>
```

TALK PRINTING ABOUT KEYS AND VALUES

Returns the matching nodes of the path expression for the given yaml file (or STDIN).

See docs for [path expression](path_expressions.md) for more details.

## Basic
Given a sample.yaml file of:
```yaml
b:
  c: 2
```
then
```bash
yq r sample.yaml b.c
```
will output the value of '2'.

## From Stdin
Given a sample.yaml file of:
```bash
cat sample.yaml | yq r - b.c
```
will output the value of '2'.


## Multiple Documents
### Reading from a single document
Given a sample.yaml file of:
```yaml
something: else
---
b:
  c: 2
```
then
```bash
yq r -d1 sample.yaml b.c
```
will output the value of '2'.

### Read from all documents
Reading all documents will return the result as an array. This can be converted to json using the '-j' flag if desired.

Given a sample.yaml file of:
```yaml
name: Fred
age: 22
---
name: Stella
age: 23
---
name: Android
age: 232
```
then
```bash
yq r -d'*' sample.yaml name
```
will output:
```
- Fred
- Stella
- Android
```