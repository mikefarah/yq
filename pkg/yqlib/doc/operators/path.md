# Path

The `path` operator can be used to get the traversal paths of matching nodes in an expression. The path is returned as an array, which if traversed in order will lead to the matching node.

You can get the key/index of matching nodes by using the `path` operator to return the path array then piping that through `.[-1]` to get the last element of that array, the key.

Use `setpath` to set a value to the path array returned by `path`, and similarly `delpaths` for an array of path arrays.


## Map path
Given a sample.yml file of:
```yaml
a:
  b: cat
```
then
```bash
yq '.a.b | path' sample.yml
```
will output
```yaml
- a
- b
```

## Get map key
Given a sample.yml file of:
```yaml
a:
  b: cat
```
then
```bash
yq '.a.b | path | .[-1]' sample.yml
```
will output
```yaml
b
```

## Array path
Given a sample.yml file of:
```yaml
a:
  - cat
  - dog
```
then
```bash
yq '.a.[] | select(. == "dog") | path' sample.yml
```
will output
```yaml
- a
- 1
```

## Get array index
Given a sample.yml file of:
```yaml
a:
  - cat
  - dog
```
then
```bash
yq '.a.[] | select(. == "dog") | path | .[-1]' sample.yml
```
will output
```yaml
1
```

## Print path and value
Given a sample.yml file of:
```yaml
a:
  - cat
  - dog
  - frog
```
then
```bash
yq '.a[] | select(. == "*og") | [{"path":path, "value":.}]' sample.yml
```
will output
```yaml
- path:
    - a
    - 1
  value: dog
- path:
    - a
    - 2
  value: frog
```

## Set path
Given a sample.yml file of:
```yaml
a:
  b: cat
```
then
```bash
yq 'setpath(["a", "b"]; "things")' sample.yml
```
will output
```yaml
a:
  b: things
```

## Set on empty document
Running
```bash
yq --null-input 'setpath(["a", "b"]; "things")'
```
will output
```yaml
a:
  b: things
```

## Set path to extract prune deep paths
Like pick but recursive. This uses `ireduce` to deeply set the selected paths into an empty object,

Given a sample.yml file of:
```yaml
parentA: bob
parentB:
  child1: i am child1
  child2: i am child2
parentC:
  child1: me child1
  child2: me child2
```
then
```bash
yq '(.parentB.child2, .parentC.child1) as $i
  ireduce({}; setpath($i | path; $i))' sample.yml
```
will output
```yaml
parentB:
  child2: i am child2
parentC:
  child1: me child1
```

## Set array path
Given a sample.yml file of:
```yaml
a:
  - cat
  - frog
```
then
```bash
yq 'setpath(["a", 0]; "things")' sample.yml
```
will output
```yaml
a:
  - things
  - frog
```

## Set array path empty
Running
```bash
yq --null-input 'setpath(["a", 0]; "things")'
```
will output
```yaml
a:
  - things
```

## Delete path
Notice delpaths takes an _array_ of paths.

Given a sample.yml file of:
```yaml
a:
  b: cat
  c: dog
  d: frog
```
then
```bash
yq 'delpaths([["a", "c"], ["a", "d"]])' sample.yml
```
will output
```yaml
a:
  b: cat
```

## Delete array path
Given a sample.yml file of:
```yaml
a:
  - cat
  - frog
```
then
```bash
yq 'delpaths([["a", 0]])' sample.yml
```
will output
```yaml
a:
  - frog
```

## Delete - wrong parameter
delpaths does not work with a single path array

Given a sample.yml file of:
```yaml
a:
  - cat
  - frog
```
then
```bash
yq 'delpaths(["a", 0])' sample.yml
```
will output
```bash
Error: DELPATHS: expected entry [0] to be a sequence, but its a !!str. Note that delpaths takes an array of path arrays, e.g. [["a", "b"]]
```

