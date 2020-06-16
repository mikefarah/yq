---
description: Merge multiple yaml files into a one
---

# Merge

Yaml files can be merged using the 'merge' command. Each additional file merged with the first file will set values for any key not existing already or where the key has no value.

```text
yq m <yaml_file> <yaml_file2> <yaml_file3>...
```

## Merge example

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
yq merge data1.yaml data2.yaml
```

will output:

```yaml
a: simple
b: [1, 2]
c:
  test: 1
```

## Updating files in-place

```bash
yq m -i data1.yaml data2.yaml
```

will update the data1.yaml file with the merged result.

## Overwrite values

Given a data1.yaml file of:

```yaml
a: simple
b: [1, 2]
d: left alone
```

and data2.yaml file of:

```yaml
a: other
b: [3, 4]
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
b: [3, 4]
c:
  test: 1
d: left alone
```

Notice that 'b' does not result in the merging of the values within an array.

## Append values with arrays

Given a data1.yaml file of:

```yaml
a: simple
b: [1, 2]
d: hi
```

and data2.yaml file of:

```yaml
a: something
b: [3, 4]
c:
  test: 2
  other: true
```

then

```bash
yq m -a data1.yaml data2.yaml
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

Note that the 'b' array has concatenated the values from the second data file. Also note that other map keys are not overridden \(field a\).

## Auto-create

By default, `yq` will automatically create any missing entries in the target yaml file. This can be turned off so that only matching paths are merged in. When turned off - you will most likely want to use the [override flag](merge.md#overwrite-values).

Given a data1.yml file of:

```yaml
a: thing
b: something else
```

and data2.yml file of:

```yaml
b: new value
d: not in original
```

Then

```yaml
yq m --overwrite --autocreate=false data1.yml data2.yml
```

will yield

```yaml
a: thing
b: new value
```



## Multiple Documents

### Merge into single document

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

### Merge into all documents

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

