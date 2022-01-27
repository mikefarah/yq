# Has

This is operation that returns true if the key exists in a map (or index in an array), false otherwise.

## Has map key
Given a sample.yml file of:
```yaml
- a: yes
- a: ~
- a:
- b: nope
```
then
```bash
yq '.[] | has("a")' sample.yml
```
will output
```yaml
true
true
true
false
```

## Select, checking for existence of deep paths
Simply pipe in parent expressions into `has`

Given a sample.yml file of:
```yaml
- a:
    b:
      c: cat
- a:
    b:
      d: dog
```
then
```bash
yq '.[] | select(.a.b | has("c"))' sample.yml
```
will output
```yaml
a:
  b:
    c: cat
```

## Has array index
Given a sample.yml file of:
```yaml
- []
- [1]
- [1, 2]
- [1, null]
- [1, 2, 3]

```
then
```bash
yq '.[] | has(1)' sample.yml
```
will output
```yaml
false
false
true
true
true
```

