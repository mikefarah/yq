The path operator can be used to get the traversal paths of matching nodes in an expression. The path is returned as an array, which if traversed in order will lead to the matching node.

You can get the key/index of matching nodes by using the `path` operator to return the path array then piping that through `.[-1]` to get the last element of that array, the key.

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

## Get map key
Given a sample.yml file of:
```yaml
a:
  b: cat
```
then
```bash
yq eval '.a.b | path | .[-1]' sample.yml
```
will output
```yaml
b
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

## Get array index
Given a sample.yml file of:
```yaml
a:
  - cat
  - dog
```
then
```bash
yq eval '.a.[] | select(. == "dog") | path | .[-1]' sample.yml
```
will output
```yaml
1
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
