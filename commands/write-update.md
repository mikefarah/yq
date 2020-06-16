---
description: >-
  Updates all the matching nodes of path expression in a yaml file to the
  supplied value.
---

# Write

```bash
yq w <yaml_file> <path_expression> <new value>
```

See docs for [path expression](../usage/path-expressions.md) and [value parsing](../usage/value-parsing.md) for more details, including controlling quotes and tags.

## Basic

Given a sample.yaml file of:

```yaml
b:
  c: 2
```

then

```bash
yq w sample.yaml b.c cat
```

will output:

```yaml
b:
  c: cat
```

### Updating files in-place

```bash
yq w -i sample.yaml b.c cat
```

will update the sample.yaml file so that the value of 'c' is cat.

## From STDIN

```bash
cat sample.yaml | yq w - b.c blah
```

## Adding new fields

Any missing fields in the path will be created on the fly.

Given a sample.yaml file of:

```yaml
b:
  c: 2
```

then

```bash
yq w sample.yaml b.d[+] "new thing"
```

will output:

```yaml
b:
  c: cat
  d:
    - new thing
```

## Appending value to an array field

Given a sample.yaml file of:

```yaml
b:
  c: 2
  d:
    - new thing
    - foo thing
```

then

```bash
yq w sample.yaml "b.d[+]" "bar thing"
```

will output:

```yaml
b:
  c: cat
  d:
    - new thing
    - foo thing
    - bar thing
```

Note that the path is in quotes to avoid the square brackets being interpreted by your shell.

## Multiple Documents

### Update a single document

Given a sample.yaml file of:

```yaml
something: else
---
b:
  c: 2
```

then

```bash
yq w -d1 sample.yaml b.c 5
```

will output:

```yaml
something: else
---
b:
  c: 5
```

### Update all documents

Given a sample.yaml file of:

```yaml
something: else
---
b:
  c: 2
```

then

```bash
yq w -d'*' sample.yaml b.c 5
```

will output:

```yaml
something: else
b:
  c: 5
---
b:
  c: 5
```

## Writing Anchors

The `---anchorName` flag can be used to set the anchor name of a node

Given a sample document of:

```yaml
commonStuff:
    flavour: vanilla
```

Then:

```bash
yq write sample.yaml commonStuff --anchorName=commonBits
```

Will yield

```yaml
commonStuff: &commonBits
    flavour: vanilla
```

## Writing Aliases

The `--makeAlias` flag can create \(or update\) a node to be an alias to an anchor.

Given a sample file of:

```yaml
commonStuff: &commonBits
    flavour: vanilla
```

Then

```bash
yq write sample.yaml --makeAnchor foo commonBits
```

Will yield:

```yaml
commonStuff: &commonBits
    flavour: vanilla
foo: *commonBits
```

## Updating only styles/tags without affecting values

You can use the write command to update the quoting style of nodes, or their tags, without re-specifying the values. This is done by omitting the value argument:

Given a sample document:

```yaml
a:
  c: things
  d: other things
```

Then

```bash
yq write sample.yaml --style=single a.*
```

Will yield:

```yaml
a:
  c: 'things'
  d: 'other things'
```

## Using a script file to update

Given a sample.yaml file of:

```yaml
b:
  d: be gone
  c: 2
  e:
    - name: Billy Bob # comment over here
```

and a script update\_instructions.yaml of:

```yaml
- command: update 
  path: b.c
  value:
    #great 
    things: frog # wow!
- command: delete
  path: b.d
```

then

```bash
yq w -s update_instructions.yaml sample.yaml
```

will output:

```yaml
b:
  c:
    #great
    things: frog # wow!
  e:
  - name: Billy Bob # comment over here
```

And, of course, you can pipe the instructions in using '-':

```bash
cat update_instructions.yaml | yq w -s - sample.yaml
```

