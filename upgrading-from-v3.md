# Upgrading from V3

Version 4 of `yq` is quite different from previous versions (and I apologise for that) - however it will be very familiar if you have used `jq` before as it now uses a similar syntax. Most commands that you could do in `v3` are longer in `v4` as a result of having a more expressive syntax language.

Note that `v4` by default now:

* prints all documents of a yaml file.
* prints in color (when outputting to a terminal).
* document separators are printed out by default

## How to do v3 things in v4:

In `v3` yq had seperate commands for reading/writing/deleting and more. In `v4` all these have been embedded into a single expression you specify to either the `eval` command (which runs the expression against each yaml document for each file given in sequence) or the `eval-all` command, which reads all documents of all files, and runs the given expression once.

Many flags from `v3` have been put into the expression language, for instance `stripComments` allowing you to specify which nodes to strip comments from instead of only being able to apply the flag to the entire document.

Lets have a look at the commands for the most common tasks:

### Navigating

v3:

```
yq r sample.yaml 'a.b.c'
```

v4:

```
yq '.a.b.c' sample.yaml
```

### Reading with default value

v3:

```
yq r sample.yaml --defaultValue frog path.not.there
```

v4: (use the [alternative](broken-reference) operator)

```
yq '.path.not.there // "frog"' sample.yaml
```



### Finding nodes

v3:

```bash
yq r sample.yaml 'a.(b.d==cat).f'
```

v4:

```bash
yq '.a | select(.b.d == "cat") | .f' sample.yaml
```

### Recursively match nodes

v3:

```
yq r sample.yaml 'thing.**.name'
```

v4:

```
yq '.thing | .. | select(has("name"))' sample.yaml
```

### Multiple documents

v3:

```bash
yq r -d1 sample.yaml 'b.c'
```

v4 (via the document index operator):

```bash
yq 'select(documentIndex == 1) | .b.c' sample.yml
```

### Updating / writing documents

v3:

```
yq w sample.yaml 'a.b.c' fred
```

v4:

```
yq '.a.b.c = "fred"' sample.yaml
```

### Deleting documents

v3:

```bash
yq d sample.yaml 'a.b.c'
```

v4:

```bash
yq 'del(.a.b.c)' sample.yaml
```

### Merging documents

Like `jq`, merge is done via the multiply operator. In yq, the merge operator can take extra options to modify how it works.

For `v3` compatability, use the `n` option to only merge in new fields.

```bash
yq '. *n load("file2.yaml")' file1.yaml
```

See the [multiply documentation](https://mikefarah.gitbook.io/yq/operators/multiply-merge) for more example and options.

### Prefix yaml

Use the [Create / Collect Into Object ](broken-reference)operator to create a new object with the desired prefix.&#x20;

v3:

```
yq p data1.yaml c.d
```

v4:

```
yq '{"c": {"d": . }}' data1.yml
```

### Create new yaml documents

Note that in v4 you can no longer run expressions against an empty file to populate it - because the file is empty, there are no matches for `yq` to run through the expression pipeline - for what it's worth, this is what `jq` does as well. Instead use the `--null-input/-n` flag and pipe out the results to the file you want directly (see example below).

v3:

```
yq n b.c cat
```

v4:

```
yq -n '.b.c = "cat"'
```

### Validate documents

v3:

```
yq validate some.file
```

v4:

```
yq 'true' some.file > /dev/null
```

Note that passing 'true' as the expression saves having to reencode the yaml (only to pipe it to stdout). In v4 you can also do a slightly more sophisticated validation and assert the tag on the root level, so you can ensure the yaml file is a map or array at the top level:

```
yq --exit-status 'tag == "!!map" or tag== "!!seq"' some.file > /dev/null
```

### Comparing yaml files

v3:

```
yq compare --prettyPrint file1.yml file2.yml 
```

v4:

In v4 there is no built in compare command, instead it relies on using diff. The downside is longer syntax, the upside is that you can use the full power of diff.

```
diff <(yq -P file1.yml) <(yq -P file2.yml)
```

### Script files

v3 had a script feature that let you run an array of commands specified in a file in one go. The format for this looked like

```yaml
- command: update 
  path: a.key1
  value: things
- command: delete
  path: a.ab.key2
```

V4 doesn't have a similar feature, however the fact that you can run multiple operations in a single expression makes it easier to come up with a shell script that does the same thing:

```bash
#!/bin/bash

yq '
  .a.key1 = "things" |
  del(.a.ab.key2)
' ./examples/data1.yaml
```

### Some new things you can do in v4:

Construct dynamic yaml [maps ](broken-reference)and [arrays ](broken-reference)based on input yaml

Using the [union ](broken-reference)operator, you can run multiple updates in one go and read multiple paths in one go

Fine grain merging of maps using the [multiply](broken-reference) operator

Read and and control yaml metadata better (e.g. [tags](broken-reference), [paths](broken-reference), [document indexes](broken-reference), [anchors and aliases](broken-reference), [comments](broken-reference)).

Work with multiple files (not just for merge)

The underlying expression language is much more powerful than `v3` so expect to see more features soon!





###
