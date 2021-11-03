# Document Index

Use the `documentIndex` operator (or the `di` shorthand) to select nodes of a particular document.

## Retrieve a document index
Given a sample.yml file of:
```yaml
a: cat
---
a: frog
```
then
```bash
yq eval '.a | documentIndex' sample.yml
```
will output
```yaml
0
---
1
```

## Retrieve a document index, shorthand
Given a sample.yml file of:
```yaml
a: cat
---
a: frog
```
then
```bash
yq eval '.a | di' sample.yml
```
will output
```yaml
0
---
1
```

## Filter by document index
Given a sample.yml file of:
```yaml
a: cat
---
a: frog
```
then
```bash
yq eval 'select(documentIndex == 1)' sample.yml
```
will output
```yaml
a: frog
```

## Filter by document index shorthand
Given a sample.yml file of:
```yaml
a: cat
---
a: frog
```
then
```bash
yq eval 'select(di == 1)' sample.yml
```
will output
```yaml
a: frog
```

## Print Document Index with matches
Given a sample.yml file of:
```yaml
a: cat
---
a: frog
```
then
```bash
yq eval '.a | ({"match": ., "doc": documentIndex})' sample.yml
```
will output
```yaml
match: cat
doc: 0
match: frog
doc: 1
```

