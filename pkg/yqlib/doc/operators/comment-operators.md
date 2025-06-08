# Comment Operators

Use these comment operators to set or retrieve comments. Note that line comments on maps/arrays are actually set on the _key_ node as opposed to the _value_ (map/array). See below for examples.

Like the `=` and `|=` assign operators, the same syntax applies when updating comments:

### plain form: `=`
This will set the LHS nodes' comments equal to the expression on the RHS. The RHS is run against the matching nodes in the pipeline

### relative form: `|=` 
This is similar to the plain form, but it evaluates the RHS with _each matching LHS node as context_. This is useful if you want to set the comments as a relative expression of the node, for instance its value or path.

## Set line comment
Set the comment on the key node for more reliability (see below).

Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq '.a line_comment="single"' sample.yml
```
will output
```yaml
a: cat # single
```

## Set line comment of a maps/arrays
For maps and arrays, you need to set the line comment on the _key_ node. This will also work for scalars.

Given a sample.yml file of:
```yaml
a:
  b: things
```
then
```bash
yq '(.a | key) line_comment="single"' sample.yml
```
will output
```yaml
a: # single
  b: things
```

## Use update assign to perform relative updates
Given a sample.yml file of:
```yaml
a: cat
b: dog
```
then
```bash
yq '.. line_comment |= .' sample.yml
```
will output
```yaml
a: cat # cat
b: dog # dog
```

## Where is the comment - map key example (legacy-v3)
The underlying yaml parser can assign comments in a document to surprising nodes. Use an expression like this to find where you comment is. 'p' indicates the path, 'isKey' is if the node is a map key (as opposed to a map value).
From this, you can see the 'hello-world-comment' is actually on the 'hello' key

Given a sample.yml file of:
```yaml
hello: # hello-world-comment
  message: world
```
then
```bash
yq '[... | {"p": path | join("."), "isKey": is_key, "hc": headComment, "lc": lineComment, "fc": footComment}]' sample.yml
```
will output
```yaml
- p: ""
  isKey: false
  hc: ""
  lc: ""
  fc: ""
- p: hello
  isKey: true
  hc: ""
  lc: hello-world-comment
  fc: ""
- p: hello
  isKey: false
  hc: ""
  lc: ""
  fc: ""
- p: hello.message
  isKey: true
  hc: ""
  lc: ""
  fc: ""
- p: hello.message
  isKey: false
  hc: ""
  lc: ""
  fc: ""
```

## Retrieve comment - map key example (legacy-v3)
From the previous example, we know that the comment is on the 'hello' _key_ as a lineComment

Given a sample.yml file of:
```yaml
hello: # hello-world-comment
  message: world
```
then
```bash
yq '.hello | key | line_comment' sample.yml
```
will output
```yaml
hello-world-comment
```

## Where is the comment - array example (legacy-v3)
The underlying yaml parser can assign comments in a document to surprising nodes. Use an expression like this to find where you comment is. 'p' indicates the path, 'isKey' is if the node is a map key (as opposed to a map value).
From this, you can see the 'under-name-comment' is actually on the first child

Given a sample.yml file of:
```yaml
name:
  - first-array-child
```
then
```bash
yq '[... | {"p": path | join("."), "isKey": is_key, "hc": headComment, "lc": lineComment, "fc": footComment}]' sample.yml
```
will output
```yaml
- p: ""
  isKey: false
  hc: ""
  lc: ""
  fc: ""
- p: name
  isKey: true
  hc: ""
  lc: ""
  fc: ""
- p: name
  isKey: false
  hc: ""
  lc: ""
  fc: ""
- p: name.0
  isKey: false
  hc: ""
  lc: ""
  fc: ""
```

## Retrieve comment - array example (legacy-v3)
From the previous example, we know that the comment is on the first child as a headComment

Given a sample.yml file of:
```yaml
name:
  - first-array-child
```
then
```bash
yq '.name[0] | headComment' sample.yml
```
will output
```yaml

```

## Where is the comment - array example (goccy)
Goccy parser has stricter comment association rules. The 'under-name-comment' is not associated with the first array child.

Given a sample.yml file of:
```yaml
name:
  - first-array-child
```
then
```bash
yq '[... | {"p": path | join("."), "isKey": is_key, "hc": headComment, "lc": lineComment, "fc": footComment}]' sample.yml
```
will output
```yaml
- p: ""
  isKey: false
  hc: ""
  lc: ""
  fc: ""
- p: name
  isKey: true
  hc: ""
  lc: ""
  fc: ""
- p: name
  isKey: false
  hc: ""
  lc: ""
  fc: ""
- p: name.0
  isKey: false
  hc: ""
  lc: ""
  fc: ""
```

## Retrieve comment - array example (goccy)
From the previous example, goccy parser does not associate the comment with the first child

Given a sample.yml file of:
```yaml
name:
  - first-array-child
```
then
```bash
yq '.name[0] | headComment' sample.yml
```
will output
```yaml

```

## Set head comment
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq '. head_comment="single"' sample.yml
```
will output
```yaml
# single
a: cat
```

## Set head comment of a map entry
Given a sample.yml file of:
```yaml
f: foo
a:
  b: cat
```
then
```bash
yq '(.a | key) head_comment="single"' sample.yml
```
will output
```yaml
f: foo
# single
a:
  b: cat
```

## Set foot comment, using an expression
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq '. foot_comment=.a' sample.yml
```
will output
```yaml
a: cat
# cat
```

## Remove (strip) all comments
Given a sample.yml file of:
```yaml
a: cat
b: dog # leave this
# footer
```
then
```bash
yq '... line_comment=""' sample.yml
```
will output
```yaml
a: cat
b: dog
# footer
```

## Remove (strip) all comments
Note the use of `...` to ensure key nodes are included.

Given a sample.yml file of:
```yaml
# hi

a: cat # comment
b: # key comment
```
then
```bash
yq '... comments=""' sample.yml
```
will output
```yaml
a: cat
b:
```

## Get line comment
Given a sample.yml file of:
```yaml
# welcome!

a: cat # meow
# have a great day
```
then
```bash
yq '.a | line_comment' sample.yml
```
will output
```yaml
meow
```

## Get head comment
Given a sample.yml file of:
```yaml
# welcome!

a: cat # meow

# have a great day
```
then
```bash
yq '. | head_comment' sample.yml
```
will output
```yaml
welcome!

```

## Head comment with document split
Given a sample.yml file of:
```yaml
# welcome!
---
# bob
a: cat # meow

# have a great day
```
then
```bash
yq 'head_comment' sample.yml
```
will output
```yaml
welcome!
bob
```

## Get foot comment
Given a sample.yml file of:
```yaml
# welcome!

a: cat # meow

# have a great day
# no really
```
then
```bash
yq '. | foot_comment' sample.yml
```
will output
```yaml

```

## Comment preservation during data operations - both parsers
Both parsers preserve structural integrity while handling comments differently

Given a sample.yml file of:
```yaml
# header
a: cat # inline
b: dog
# footer
```
then
```bash
yq '.c = "new"' sample.yml
```
will output
```yaml
# header
a: cat # inline
b: dog
# footer

c: new
```

## Comment before array items - legacy-v3 behavior
legacy-v3 associates comments that precede array elements

Given a sample.yml file of:
```yaml
items:
  - name: first
    value: 100
```
then
```bash
yq '.items[0] | head_comment' sample.yml
```
will output
```yaml

```

## Comment before array items - goccy behavior
Goccy has stricter rules and does not associate this comment with the array element

Given a sample.yml file of:
```yaml
items:
  - name: first
    value: 100
```
then
```bash
yq '.items[0] | head_comment' sample.yml
```
will output
```yaml

```

## Comment between map keys - both parsers
Both parsers handle comments between sibling elements consistently

Given a sample.yml file of:
```yaml
key1: value1
key2: value2
```
then
```bash
yq '.key2 | head_comment' sample.yml
```
will output
```yaml

```

## Complex comment scenario - legacy-v3
Complex document with multiple comment types - legacy-v3 behavior

Given a sample.yml file of:
```yaml
# Document header
config:
  - name: service1
    port: 8080
  - name: service2
    port: 9090

# Document footer
```
then
```bash
yq '.config[0].port | head_comment' sample.yml
```
will output
```yaml

```

## Complex comment scenario - goccy
Complex document with multiple comment types - goccy behavior

Given a sample.yml file of:
```yaml
# Document header
config:
  - name: service1
    port: 8080
  - name: service2
    port: 9090

# Document footer
```
then
```bash
yq '.config[0].port | head_comment' sample.yml
```
will output
```yaml

```

