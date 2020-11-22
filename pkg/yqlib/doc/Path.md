The path operator can be used to get the traversal paths of matching nodes in an expression. The path is returned as an array, which if traversed in order will lead to the matching node.
## Map path
Given a sample.yml file of:
```yaml
a:
  b: cat
```
then
```bash
yq eval '.a.b | path' sample.yml
```
will output
```yaml
- a
- b
```

## Array path
Given a sample.yml file of:
```yaml
a:
  - cat
  - dog
```
then
```bash
yq eval '.a.[] | select(. == "dog") | path' sample.yml
```
will output
```yaml
- a
- 1
```

## Print path and value
Given a sample.yml file of:
```yaml
a:
  - cat
  - dog
  - frog
```
then
```bash
yq eval '.a.[] | select(. == "*og") | [{"path":path, "value":.}]' sample.yml
```
will output
```yaml
- path:
    - a
    - 1
  value: dog
- path:
    - a
    - 2
  value: frog
```

