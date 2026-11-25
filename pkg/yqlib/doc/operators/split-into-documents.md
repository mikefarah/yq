# Split into Documents

This operator splits all matches into separate documents

## Split after assigning a variable
Given a sample.yml file of:
```yaml
- 1
- 2
```
then
```bash
yq '.[] | . as $_ | split_doc' sample.yml
```
will output
```yaml
1
---
2
```

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

