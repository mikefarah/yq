```
yq d <yaml_file|json_file> <path_to_delete>
```
{!snippets/works_with_json.md!}

### To Stdout
Given a sample.yaml file of:
```yaml
b:
  c: 2
  apples: green
```
then
```bash
yq d sample.yaml b.c
```
will output:
```yaml
b:
  apples: green
```

### From STDIN
```bash
cat sample.yaml | yq d - b.c
```

### Deleting array elements
Given a sample.yaml file of:
```yaml
b:
  c: 
    - 1
    - 2
    - 3
```
then
```bash
yq d sample.yaml 'b.c[1]'
```
will output:
```yaml
b:
  c:
  - 1
  - 3
```

### Deleting nodes in-place
Given a sample.yaml file of:
```yaml
b:
  c: 2
  apples: green
```
then
```bash
yq d -i sample.yaml b.c
```
will update the sample.yaml file so that the 'c' node is deleted


{!snippets/keys_with_dots.md!}
