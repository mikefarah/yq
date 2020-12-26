Use the `documentIndex` operator to select nodes of a particular document.
## Retrieve a document index
Given a sample.yml file of:
```yaml
a: cat
'': null
---
a: frog
'': null
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

## Filter by document index
Given a sample.yml file of:
```yaml
a: cat
'': null
---
a: frog
'': null
```
then
```bash
yq eval 'select(. | documentIndex == 1)' sample.yml
```
will output
```yaml
a: frog
'': null
```

## Print Document Index with matches
Given a sample.yml file of:
```yaml
a: cat
'': null
---
a: frog
'': null
```
then
```bash
yq eval '.a | ({"match": ., "doc": (. | documentIndex)})' sample.yml
```
will output
```yaml
'': null
'': null
```

