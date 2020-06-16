# Create

 Yaml files can be created using the 'new' command. This works in the same way as the write command, but you don't pass in an existing Yaml file. Currently this does not support creating multiple documents in a single yaml file.

```text
yq n <path> <new value>
```

### Creating a simple yaml file[¶](create.md#creating-a-simple-yaml-file) <a id="creating-a-simple-yaml-file"></a>

```text
yq n b.c cat
```

will output:

```text
b:
  c: cat
```

### Creating using a create script[¶](create.md#creating-using-a-create-script) <a id="creating-using-a-create-script"></a>

Create scripts follow the same format as the update scripts.

Given a script create\_instructions.yaml of:

```text
b.c: 3
b.e[+].name: Howdy Partner
```

then

```text
yq n -s create_instructions.yaml
```

will output:

```text
b:
  c: 3
  e:
    - name: Howdy Partner
```

You can also pipe the instructions in:

```text
cat create_instructions.yaml | yq n -s -
```

### Keys with dots[¶](create.md#keys-with-dots) <a id="keys-with-dots"></a>

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

### Keys \(and values\) with leading dashes[¶](create.md#keys-and-values-with-leading-dashes) <a id="keys-and-values-with-leading-dashes"></a>

If a key or value has leading dashes, yq won't know that you are passing a value as opposed to a flag \(and you will get a 'bad flag syntax' error\).

To fix that, you will need to tell it to stop processing flags by adding '--' after the last flag like so:

```text
yq n -t -- --key --value
```

Will result in

``` --key: --value``

