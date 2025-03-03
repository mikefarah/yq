# Entries

Similar to the same named functions in `jq` these functions convert to/from an object and an array of key-value pairs. This is most useful for performing operations on keys of maps.

Use `with_entries(op)` as a syntatic sugar for doing `to_entries | op | from_entries`.

## to_entries Map
Given a sample.yml file of:
```yaml
a: 1
b: 2
```
then
```bash
yq 'to_entries' sample.yml
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
yq 'to_entries' sample.yml
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
yq 'to_entries' sample.yml
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
yq 'to_entries | from_entries' sample.yml
```
will output
```yaml
a: 1
b: 2
```

## from_entries with numeric key indices
from_entries always creates a map, even for numeric keys

Given a sample.yml file of:
```yaml
- a
- b
```
then
```bash
yq 'to_entries | from_entries' sample.yml
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
yq 'with_entries(.key |= "KEY_" + .)' sample.yml
```
will output
```yaml
KEY_a: 1
KEY_b: 2
```

## Use with_entries to update keys recursively
We use (.. | select(tag="map")) to find all the maps in the doc, then |= to update each one of those maps. In the update, with_entries is used.

Given a sample.yml file of:
```yaml
a: 1
b:
  b_a: nested
  b_b: thing
```
then
```bash
yq '(.. | select(tag=="!!map")) |= with_entries(.key |= "KEY_" + .)' sample.yml
```
will output
```yaml
KEY_a: 1
KEY_b:
  KEY_b_a: nested
  KEY_b_b: thing
```

## Custom sort map keys
Use to_entries to convert to an array of key/value pairs, sort the array using sort/sort_by/etc, and convert it back.

Given a sample.yml file of:
```yaml
a: 1
c: 3
b: 2
```
then
```bash
yq 'to_entries | sort_by(.key) | reverse | from_entries' sample.yml
```
will output
```yaml
c: 3
b: 2
a: 1
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
yq 'with_entries(select(.value | has("b")))' sample.yml
```
will output
```yaml
a:
  b: bird
```

