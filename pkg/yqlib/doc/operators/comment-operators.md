# Comment Operators

Use these comment operators to set or retrieve comments. Note that line comments on maps/arrays are actually set on the _key_ node as opposed to the _value_ (map/array). See below for examples.

Like the `=` and `|=` assign operators, the same syntax applies when updating comments:

### plain form: `=`
This will set the LHS nodes' comments equal to the expression on the RHS. The RHS is run against the matching nodes in the pipeline

### relative form: `|=` 
This is similar to the plain form, but it evaluates the RHS with _each matching LHS node as context_. This is useful if you want to set the comments as a relative expression of the node, for instance its value or path.

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

