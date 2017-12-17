Yaml files can be merged using the 'merge' command. Each additional file merged with the first file will
set values for any key not existing already or where the key has no value.

```
yq m <yaml_file|json_file> <path>...
```
{!snippets/works_with_json.md!}

### To Stdout
Given a data1.yaml file of:
```yaml
a: simple
b: [1, 2]
```
and data2.yaml file of:
```yaml
a: other
c:
  test: 1
```
then
```bash
yq m data1.yaml data2.yaml
```
will output:
```yaml
a: simple
b: [1, 2]
c:
  test: 1
```

### Updating files in-place
Given a data1.yaml file of:
```yaml
a: simple
b: [1, 2]
```
and data2.yaml file of:
```yaml
a: other
c:
  test: 1
```
then
```bash
yq m -i data1.yaml data2.yaml
```
will update the data1.yaml file so that the value of 'c' is 'test: 1'.

### Overwrite values
Given a data1.yaml file of:
```yaml
a: simple
b: [1, 2]
```
and data2.yaml file of:
```yaml
a: other
c:
  test: 1
```
then
```bash
yq m -x data1.yaml data2.yaml
```
will output:
```yaml
a: other
b: [1, 2]
c:
  test: 1
```

### Overwrite values with arrays
Given a data1.yaml file of:
```yaml
a: simple
b: [1, 2]
```
and data3.yaml file of:
```yaml
b: [2, 3, 4]
c:
  test: 2
  other: true
d: false
```
then
```bash
yq m -x data1.yaml data3.yaml
```
will output:
```yaml
a: simple
b: [2, 3, 4]
c:
  test: 2
  other: true
d: false
```

Notice that 'b' does not result in the merging of the values within an array. The underlying library does not
currently handle merging values within an array.
