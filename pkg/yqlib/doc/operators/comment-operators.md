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
a:
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

## Where is the comment - map key example
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
  isKey: null
  true: null
  hc: null
  "": null
  lc: null
  hello-world-comment: null
  fc: null
- p: hello
  isKey: false
  hc: ""
  lc: ""
  fc: ""
- p: hello.message
  isKey: null
  true: null
  hc: null
  "": null
  lc: null
  fc: null
- p: hello.message
  isKey: false
  hc: ""
  lc: ""
  fc: ""
```

## Retrieve comment - map key example
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

## Where is the comment - array example
The underlying yaml parser can assign comments in a document to surprising nodes. Use an expression like this to find where you comment is. 'p' indicates the path, 'isKey' is if the node is a map key (as opposed to a map value).
From this, you can see the 'under-name-comment' is actually on the first child

Given a sample.yml file of:
```yaml
name:
  # under-name-comment
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
  isKey: null
  true: null
  hc: null
  "": null
  lc: null
  fc: null
- p: name
  isKey: false
  hc: ""
  lc: ""
  fc: ""
- p: name.0
  isKey: false
  hc: under-name-comment
  lc: ""
  fc: ""
```

## Retrieve comment - array example
From the previous example, we know that the comment is on the first child as a headComment

Given a sample.yml file of:
```yaml
name:
  # under-name-comment
  - first-array-child
```
then
```bash
yq '.name[0] | headComment' sample.yml
```
will output
```yaml
under-name-comment
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

## Remove comment
Given a sample.yml file of:
```yaml
a: cat # comment
b: dog # leave this
```
then
```bash
yq '.a line_comment=""' sample.yml
```
will output
```yaml
a: cat
b: dog # leave this
```

## Remove (strip) all comments
Note the use of `...` to ensure key nodes are included.

Given a sample.yml file of:
```yaml
# hi

a: cat # comment
# great
b: # key comment
```
then
```bash
yq '... comments=""' sample.yml
```
will output
```yaml
# hi

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
have a great day
no really
```

