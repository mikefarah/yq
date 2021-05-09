
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

