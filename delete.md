# Delete

```text
yq d <yaml_file> <path_to_delete>
```

### To Stdout[¶](delete.md#to-stdout) <a id="to-stdout"></a>

Given a sample.yaml file of:

```text
b:
  c: 2
  apples: green
```

then

```text
yq d sample.yaml b.c
```

will output:

```text
b:
  apples: green
```

### From STDIN[¶](delete.md#from-stdin) <a id="from-stdin"></a>

```text
cat sample.yaml | yq d - b.c
```

### Deleting array elements[¶](delete.md#deleting-array-elements) <a id="deleting-array-elements"></a>

Given a sample.yaml file of:

```text
b:
  c: 
    - 1
    - 2
    - 3
```

then

```text
yq d sample.yaml 'b.c[1]'
```

will output:

```text
b:
  c:
  - 1
  - 3
```

### Deleting nodes in-place[¶](delete.md#deleting-nodes-in-place) <a id="deleting-nodes-in-place"></a>

Given a sample.yaml file of:

```text
b:
  c: 2
  apples: green
```

then

```text
yq d -i sample.yaml b.c
```

will update the sample.yaml file so that the 'c' node is deleted

### Splat[¶](delete.md#splat) <a id="splat"></a>

Given a sample.yaml file of:

```text
---
bob:
  item1:
    cats: bananas
    dogs: woof
  item2:
    cats: apples
    dogs: woof2
  thing:
    cats: oranges
    dogs: woof3
```

then

```text
yq d sample.yaml bob.*.cats
```

will output:

```text
---
bob:
  item1:
    dogs: woof
  item2:
    dogs: woof2
  thing:
    dogs: woof3
```

### Prefix Splat[¶](delete.md#prefix-splat) <a id="prefix-splat"></a>

Given a sample.yaml file of:

```text
---
bob:
  item1:
    cats: bananas
    dogs: woof
  item2:
    cats: apples
    dogs: woof2
  thing:
    cats: oranges
    dogs: woof3
```

then

```text
yq d sample.yaml bob.item*.cats
```

will output:

```text
---
bob:
  item1:
    dogs: woof
  item2:
    dogs: woof2
  thing:
    cats: oranges
    dogs: woof3
```

### Array Splat[¶](delete.md#array-splat) <a id="array-splat"></a>

Given a sample.yaml file of:

```text
---
bob:
- cats: bananas
  dogs: woof
- cats: apples
  dogs: woof2
- cats: oranges
  dogs: woof3
```

then

```text
yq d sample.yaml bob.[*].cats
```

will output:

```text
---
bob:
- dogs: woof
- dogs: woof2
- dogs: woof3
```

### Multiple Documents - delete from single document[¶](delete.md#multiple-documents-delete-from-single-document) <a id="multiple-documents-delete-from-single-document"></a>

Given a sample.yaml file of:

```text
something: else
field: leaveMe
---
b:
  c: 2
field: deleteMe
```

then

```text
yq w -d1 sample.yaml field
```

will output:

```text
something: else
field: leaveMe
---
b:
  c: 2
```

### Multiple Documents - delete from all documents[¶](delete.md#multiple-documents-delete-from-all-documents) <a id="multiple-documents-delete-from-all-documents"></a>

Given a sample.yaml file of:

```text
something: else
field: deleteMe
---
b:
  c: 2
field: deleteMeToo
```

then

```text
yq w -d'*' sample.yaml field
```

will output:

```text
something: else
---
b:
  c: 2
```

Note that '\*' is in quotes to avoid being interpreted by your shell.

### Keys with dots[¶](delete.md#keys-with-dots) <a id="keys-with-dots"></a>

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

### Keys \(and values\) with leading dashes[¶](delete.md#keys-and-values-with-leading-dashes) <a id="keys-and-values-with-leading-dashes"></a>

If a key or value has leading dashes, yq won't know that you are passing a value as opposed to a flag \(and you will get a 'bad flag syntax' error\).

To fix that, you will need to tell it to stop processing flags by adding '--' after the last flag like so:

```text
yq n -t -- --key --value
```

Will result in

``` --key: --value``

