# Write

```text
yq w <yaml_file> <path> <new value>
```

### To Stdout[¶](write-update.md#to-stdout) <a id="to-stdout"></a>

Given a sample.yaml file of:

```text
b:
  c: 2
```

then

```text
yq w sample.yaml b.c cat
```

will output:

```text
b:
  c: cat
```

### From STDIN[¶](write-update.md#from-stdin) <a id="from-stdin"></a>

```text
cat sample.yaml | yq w - b.c blah
```

### Adding new fields[¶](write-update.md#adding-new-fields) <a id="adding-new-fields"></a>

Any missing fields in the path will be created on the fly.

Given a sample.yaml file of:

```text
b:
  c: 2
```

then

```text
yq w sample.yaml b.d[+] "new thing"
```

will output:

```text
b:
  c: cat
  d:
    - new thing
```

### Splat[¶](write-update.md#splat) <a id="splat"></a>

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
yq w sample.yaml bob.*.cats meow
```

will output:

```text
---
bob:
  item1:
    cats: meow
  item2:
    cats: meow
  thing:
    cats: meow
```

### Prefix Splat[¶](write-update.md#prefix-splat) <a id="prefix-splat"></a>

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
yq w sample.yaml bob.item*.cats meow
```

will output:

```text
---
bob:
  item1:
    cats: meow
  item2:
    cats: meow
  thing:
    cats: oranges
```

### Array Splat[¶](write-update.md#array-splat) <a id="array-splat"></a>

Given a sample.yaml file of:

```text
---
bob:
- cats: bananas
- cats: apples
- cats: oranges
```

then

```text
yq w sample.yaml bob[*].cats meow
```

will output:

```text
---
bob:
- cats: meow
- cats: meow
- cats: meow
```

### Appending value to an array field[¶](write-update.md#appending-value-to-an-array-field) <a id="appending-value-to-an-array-field"></a>

Given a sample.yaml file of:

```text
b:
  c: 2
  d:
    - new thing
    - foo thing
```

then

```text
yq w sample.yaml "b.d[+]" "bar thing"
```

will output:

```text
b:
  c: cat
  d:
    - new thing
    - foo thing
    - bar thing
```

Note that the path is in quotes to avoid the square brackets being interpreted by your shell.

### Multiple Documents - update a single document[¶](write-update.md#multiple-documents-update-a-single-document) <a id="multiple-documents-update-a-single-document"></a>

Given a sample.yaml file of:

```text
something: else
---
b:
  c: 2
```

then

```text
yq w -d1 sample.yaml b.c 5
```

will output:

```text
something: else
---
b:
  c: 5
```

### Multiple Documents - update all documents[¶](write-update.md#multiple-documents-update-all-documents) <a id="multiple-documents-update-all-documents"></a>

Given a sample.yaml file of:

```text
something: else
---
b:
  c: 2
```

then

```text
yq w -d'*' sample.yaml b.c 5
```

will output:

```text
something: else
b:
  c: 5
---
b:
  c: 5
```

Note that '\*' is in quotes to avoid being interpreted by your shell.

### Updating files in-place[¶](write-update.md#updating-files-in-place) <a id="updating-files-in-place"></a>

Given a sample.yaml file of:

```text
b:
  c: 2
```

then

```text
yq w -i sample.yaml b.c cat
```

will update the sample.yaml file so that the value of 'c' is cat.

### Updating multiple values with a script[¶](write-update.md#updating-multiple-values-with-a-script) <a id="updating-multiple-values-with-a-script"></a>

Given a sample.yaml file of:

```text
b:
  c: 2
  e:
    - name: Billy Bob
```

and a script update\_instructions.yaml of:

```text
b.c: 3
b.e[+].name: Howdy Partner
```

then

```text
yq w -s update_instructions.yaml sample.yaml
```

will output:

```text
b:
  c: 3
  e:
    - name: Howdy Partner
```

And, of course, you can pipe the instructions in using '-':

```text
cat update_instructions.yaml | yq w -s - sample.yaml
```

### Values starting with a hyphen \(or dash\)[¶](write-update.md#values-starting-with-a-hyphen-or-dash) <a id="values-starting-with-a-hyphen-or-dash"></a>

The flag terminator needs to be used to stop the app from attempting to parse the subsequent arguments as flags:

```text
yq w -- my.path -3
```

will output

```text
my:
  path: -3
```

### Keys with dots[¶](write-update.md#keys-with-dots) <a id="keys-with-dots"></a>

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

### Keys \(and values\) with leading dashes[¶](write-update.md#keys-and-values-with-leading-dashes) <a id="keys-and-values-with-leading-dashes"></a>

If a key or value has leading dashes, yq won't know that you are passing a value as opposed to a flag \(and you will get a 'bad flag syntax' error\).

To fix that, you will need to tell it to stop processing flags by adding '--' after the last flag like so:

```text
yq n -t -- --key --value
```

Will result in

``` --key: --value``

