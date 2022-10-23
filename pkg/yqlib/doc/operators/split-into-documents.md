# Split into Documents

This operator splits all matches into separate documents

## Split empty
Running
```bash
yq --null-input 'split_doc'
```
will output
```yaml

```

## Split array
Given a sample.yml file of:
```yaml
- a: cat
- b: dog
```
then
```bash
yq '.[] | split_doc' sample.yml
```
will output
```yaml
a: cat
---
b: dog
```

