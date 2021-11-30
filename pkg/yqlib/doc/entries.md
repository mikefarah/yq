# Entries

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

## to_entries null
Given a sample.yml file of:
```yaml
null
```
then
```bash
yq eval 'to_entries' sample.yml
```
will output
```yaml
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
And another sample another.yml file of:
```yaml
c: 1
d: 2
```
then
```bash
yq eval-all 'with_entries(.key |= "KEY_" + .)' sample.yml another.yml
```
will output
```yaml
KEY_a: 1
KEY_b: 2
KEY_c: 1
KEY_d: 2
```

## Use with_entries to filter the map
Given a sample.yml file of:
```yaml
a:
  b: bird
c:
  d: dog
```
then
```bash
yq eval 'with_entries(select(.value | has("b")))' sample.yml
```
will output
```yaml
a:
  b: bird
```

