# Prefix

 Paths can be prefixed using the 'prefix' command. The complete yaml content will be nested inside the new prefix path.

```text
yq p <yaml_file> <path>
```

### To Stdout[¶](prefix.md#to-stdout) <a id="to-stdout"></a>

Given a data1.yaml file of:

```text
a: simple
b: [1, 2]
```

then

```text
yq p data1.yaml c
```

will output:

```text
c:
  a: simple
  b: [1, 2]
```

### Arbitrary depth[¶](prefix.md#arbitrary-depth) <a id="arbitrary-depth"></a>

Given a data1.yaml file of:

```text
a:
  b: [1, 2]
```

then

```text
yq p data1.yaml c.d
```

will output:

```text
c:
  d:
    a:
      b: [1, 2]
```

### Updating files in-place[¶](prefix.md#updating-files-in-place) <a id="updating-files-in-place"></a>

Given a data1.yaml file of:

```text
a: simple
b: [1, 2]
```

then

```text
yq p -i data1.yaml c
```

will update the data1.yaml file so that the path 'c' is prefixed to all other paths.

### Multiple Documents - prefix a single document[¶](prefix.md#multiple-documents-prefix-a-single-document) <a id="multiple-documents-prefix-a-single-document"></a>

Given a data1.yaml file of:

```text
something: else
---
a: simple
b: cat
```

then

```text
yq p -d1 data1.yaml c
```

will output:

```text
something: else
---
c:
  a: simple
  b: cat
```

### Multiple Documents - prefix all documents[¶](prefix.md#multiple-documents-prefix-all-documents) <a id="multiple-documents-prefix-all-documents"></a>

Given a data1.yaml file of:

```text
something: else
---
a: simple
b: cat
```

then

```text
yq p -d'*' data1.yaml c
```

will output:

```text
c:
  something: else
---
c:
  a: simple
  b: cat
```

