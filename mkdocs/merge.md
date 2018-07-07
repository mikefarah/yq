Yaml files can be merged using the 'merge' command. Each additional file merged with the first file will
set values for any key not existing already or where the key has no value.

```
yq m <yaml_file> <path>...
```


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
b: [3, 4]
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
b: [3, 4]
c:
  test: 2
  other: true
d: false
```

Notice that 'b' does not result in the merging of the values within an array. 

### Append values with arrays
Given a data1.yaml file of:
```yaml
a: simple
b: [1, 2]
d: hi
```
and data3.yaml file of:
```yaml
a: something
b: [3, 4]
c:
  test: 2
  other: true
```
then
```bash
yq m -a data1.yaml data3.yaml
```
will output:
```yaml
a: simple
b: [1, 2, 3, 4]
c:
  test: 2
  other: true
d: hi
```

Note that the 'b' array has concatenated the values from the second data file. Also note that other map keys are not overridden (field a).

Append cannot be used with overwrite, if both flags are given then append is ignored.

### Multiple Documents - merge into single document
Currently yq only has multi-document support for the _first_ document being merged into. The remaining yaml files will have their first document selected.

Given a data1.yaml file of:
```yaml
something: else
---
a: simple
b: cat
```
and data3.yaml file of:
```yaml
b: dog
```
then
```bash
yq m -x -d1 data1.yaml data3.yaml
```
will output:
```yaml
something: else
---
a: simple
b: dog
```

### Multiple Documents - merge into all documents
Currently yq only has multi-document support for the _first_ document being merged into. The remaining yaml files will have their first document selected.

Given a data1.yaml file of:
```yaml
something: else
---
a: simple
b: cat
```
and data3.yaml file of:
```yaml
b: dog
```
then
```bash
yq m -x -d'*' data1.yaml data3.yaml
```
will output:
```yaml
b: dog
something: else
---
a: simple
b: dog
```