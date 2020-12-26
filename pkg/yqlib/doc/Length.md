Returns the lengths of the nodes. Length is defined according to the type of the node.

## String length
returns length of string

Given a sample.yml file of:
```yaml
a: cat
'': null
```
then
```bash
yq eval '.a | length' sample.yml
```
will output
```yaml
3
```

## Map length
returns number of entries

Given a sample.yml file of:
```yaml
a: cat
c: dog
'': null
```
then
```bash
yq eval 'length' sample.yml
```
will output
```yaml
3
```

## Array length
returns number of elements

Given a sample.yml file of:
```yaml
- 2
- 4
- 6
- 8
```
then
```bash
yq eval 'length' sample.yml
```
will output
```yaml
4
```

