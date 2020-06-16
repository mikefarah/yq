# Merge

 Yaml files can be merged using the 'merge' command. Each additional file merged with the first file will set values for any key not existing already or where the key has no value.

```text
yq m <yaml_file> <path>...
```

### To Stdout[¶](merge.md#to-stdout) <a id="to-stdout"></a>

Given a data1.yaml file of:

```text
a: simple
b: [1, 2]
```

and data2.yaml file of:

```text
a: other
c:
  test: 1
```

then

```text
yq m data1.yaml data2.yaml
```

will output:

```text
a: simple
b: [1, 2]
c:
  test: 1
```

### Updating files in-place[¶](merge.md#updating-files-in-place) <a id="updating-files-in-place"></a>

Given a data1.yaml file of:

```text
a: simple
b: [1, 2]
```

and data2.yaml file of:

```text
a: other
c:
  test: 1
```

then

```text
yq m -i data1.yaml data2.yaml
```

will update the data1.yaml file so that the value of 'c' is 'test: 1'.

### Overwrite values[¶](merge.md#overwrite-values) <a id="overwrite-values"></a>

Given a data1.yaml file of:

```text
a: simple
b: [1, 2]
```

and data2.yaml file of:

```text
a: other
c:
  test: 1
```

then

```text
yq m -x data1.yaml data2.yaml
```

will output:

```text
a: other
b: [1, 2]
c:
  test: 1
```

### Overwrite values with arrays[¶](merge.md#overwrite-values-with-arrays) <a id="overwrite-values-with-arrays"></a>

Given a data1.yaml file of:

```text
a: simple
b: [1, 2]
```

and data3.yaml file of:

```text
b: [3, 4]
c:
  test: 2
  other: true
d: false
```

then

```text
yq m -x data1.yaml data3.yaml
```

will output:

```text
a: simple
b: [3, 4]
c:
  test: 2
  other: true
d: false
```

Notice that 'b' does not result in the merging of the values within an array.

### Append values with arrays[¶](merge.md#append-values-with-arrays) <a id="append-values-with-arrays"></a>

Given a data1.yaml file of:

```text
a: simple
b: [1, 2]
d: hi
```

and data3.yaml file of:

```text
a: something
b: [3, 4]
c:
  test: 2
  other: true
```

then

```text
yq m -a data1.yaml data3.yaml
```

will output:

```text
a: simple
b: [1, 2, 3, 4]
c:
  test: 2
  other: true
d: hi
```

Note that the 'b' array has concatenated the values from the second data file. Also note that other map keys are not overridden \(field a\).

Append cannot be used with overwrite, if both flags are given then append is ignored.

### Multiple Documents - merge into single document[¶](merge.md#multiple-documents-merge-into-single-document) <a id="multiple-documents-merge-into-single-document"></a>

Currently yq only has multi-document support for the _first_ document being merged into. The remaining yaml files will have their first document selected.

Given a data1.yaml file of:

```text
something: else
---
a: simple
b: cat
```

and data3.yaml file of:

```text
b: dog
```

then

```text
yq m -x -d1 data1.yaml data3.yaml
```

will output:

```text
something: else
---
a: simple
b: dog
```

### Multiple Documents - merge into all documents[¶](merge.md#multiple-documents-merge-into-all-documents) <a id="multiple-documents-merge-into-all-documents"></a>

Currently yq only has multi-document support for the _first_ document being merged into. The remaining yaml files will have their first document selected.

Given a data1.yaml file of:

```text
something: else
---
a: simple
b: cat
```

and data3.yaml file of:

```text
b: dog
```

then

```text
yq m -x -d'*' data1.yaml data3.yaml
```

will output:

```text
b: dog
something: else
---
a: simple
b: dog
```

