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
yq eval '.[] | has("a")' sample.yml
```
will output
```yaml
true
true
true
false
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
yq eval '.[] | has(1)' sample.yml
```
will output
```yaml
false
false
true
true
true
```

