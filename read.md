# Read

```text
yq r <yaml_file|json_file> <path>
```

This command can take a json file as input too, and will output yaml unless specified to export as json \(-j\)

### Basic[¶](read.md#basic) <a id="basic"></a>

Given a sample.yaml file of:

```text
b:
  c: 2
```

then

```text
yq r sample.yaml b.c
```

will output the value of '2'.

### From Stdin[¶](read.md#from-stdin) <a id="from-stdin"></a>

Given a sample.yaml file of:

```text
cat sample.yaml | yq r - b.c
```

will output the value of '2'.

### Splat[¶](read.md#splat) <a id="splat"></a>

Given a sample.yaml file of:

```text
---
bob:
  item1:
    cats: bananas
  item2:
    cats: apples
  thing:
    cats: oranges
```

then

```text
yq r sample.yaml bob.*.cats
```

will output

```text
- bananas
- apples
- oranges
```

### Prefix Splat[¶](read.md#prefix-splat) <a id="prefix-splat"></a>

Given a sample.yaml file of:

```text
---
bob:
  item1:
    cats: bananas
  item2:
    cats: apples
  thing:
    cats: oranges
```

then

```text
yq r sample.yaml bob.item*.cats
```

will output

```text
- bananas
- apples
```

### Multiple Documents - specify a single document[¶](read.md#multiple-documents-specify-a-single-document) <a id="multiple-documents-specify-a-single-document"></a>

Given a sample.yaml file of:

```text
something: else
---
b:
  c: 2
```

then

```text
yq r -d1 sample.yaml b.c
```

will output the value of '2'.

### Multiple Documents - read all documents[¶](read.md#multiple-documents-read-all-documents) <a id="multiple-documents-read-all-documents"></a>

Reading all documents will return the result as an array. This can be converted to json using the '-j' flag if desired.

Given a sample.yaml file of:

```text
name: Fred
age: 22
---
name: Stella
age: 23
---
name: Android
age: 232
```

then

```text
yq r -d'*' sample.yaml name
```

will output:

```text
- Fred
- Stella
- Android
```

### Arrays[¶](read.md#arrays) <a id="arrays"></a>

You can give an index to access a specific element: e.g.: given a sample file of

```text
b:
  e:
    - name: fred
      value: 3
    - name: sam
      value: 4
```

then

```text
yq r sample.yaml 'b.e[1].name'
```

will output 'sam'

Note that the path is in quotes to avoid the square brackets being interpreted by your shell.

### Array Splat[¶](read.md#array-splat) <a id="array-splat"></a>

e.g.: given a sample file of

```text
b:
  e:
    - name: fred
      value: 3
    - name: sam
      value: 4
```

then

```text
yq r sample.yaml 'b.e[*].name'
```

will output:

```text
- fred
- sam
```

Note that the path is in quotes to avoid the square brackets being interpreted by your shell.

### Keys with dots[¶](read.md#keys-with-dots) <a id="keys-with-dots"></a>

When specifying a key that has a dot use key lookup indicator.

```text
b:
  foo.bar: 7
```

```text
yaml r sample.yaml 'b[foo.bar]'
```

```text
yaml w sample.yaml 'b[foo.bar]' 9
```

Any valid yaml key can be specified as part of a key lookup.

Note that the path is in quotes to avoid the square brackets being interpreted by your shell.

### Keys \(and values\) with leading dashes[¶](read.md#keys-and-values-with-leading-dashes) <a id="keys-and-values-with-leading-dashes"></a>

If a key or value has leading dashes, yq won't know that you are passing a value as opposed to a flag \(and you will get a 'bad flag syntax' error\).

To fix that, you will need to tell it to stop processing flags by adding '--' after the last flag like so:

```text
yq n -t -- --key --value
```

Will result in

``` --key: --value``

