Similar to the same named functions in `jq` these functions convert to/from an object and an array of key-value pairs. This is most useful for performing operations on keys of maps.
## to_entries Map
Given a sample.yml file of:
```yaml
a: 1
b: 2
```
then
```bash
yq eval 'to_entries' sample.yml
```
will output
```yaml
- key: a
  value: 1
- key: b
  value: 2
```

## to_entries Array
Given a sample.yml file of:
```yaml
- a
- b
```
then
```bash
yq eval 'to_entries' sample.yml
```
will output
```yaml
- key: 0
  value: a
- key: 1
  value: b
```

## from_entries map
Given a sample.yml file of:
```yaml
a: 1
b: 2
```
then
```bash
yq eval 'to_entries | from_entries' sample.yml
```
will output
```yaml
a: 1
b: 2
```

## from_entries with numeric key indexes
from_entries always creates a map, even for numeric keys

Given a sample.yml file of:
```yaml
- a
- b
```
then
```bash
yq eval 'to_entries | from_entries' sample.yml
```
will output
```yaml
0: a
1: b
```

## Use with_entries to update keys
Given a sample.yml file of:
```yaml
a: 1
b: 2
```
then
```bash
yq eval 'with_entries(.key |= "KEY_" + .)' sample.yml
```
will output
```yaml
KEY_a: 1
KEY_b: 2
```

