Like the multiple operator in `jq`, depending on the operands, this multiply operator will do different things. Currently only objects are supported, which have the effect of merging the RHS into the LHS.

Upcoming versions of `yq` will add support for other types of multiplication (numbers, strings).

To concatenate when merging objects, use the `*+` form (see examples below). This will recursively merge objects, appending arrays when it encounters them.

Note that when merging objects, this operator returns the merged object (not the parent). This will be clearer in the examples below.

## Merging files
Note the use of eval-all to ensure all documents are loaded into memory.

```bash
yq eval-all 'select(fileIndex == 0) * select(fileIndex == 1)' file1.yaml file2.yaml
```

## Merge objects together, returning merged result only
Given a sample.yml file of:
```yaml
a:
  field: me
  fieldA: cat
b:
  field:
    g: wizz
  fieldB: dog
```
then
```bash
yq eval '.a * .b' sample.yml
```
will output
```yaml
field:
  g: wizz
fieldA: cat
fieldB: dog
```

## Merge objects together, returning parent object
Given a sample.yml file of:
```yaml
a:
  field: me
  fieldA: cat
b:
  field:
    g: wizz
  fieldB: dog
```
then
```bash
yq eval '. * {"a":.b}' sample.yml
```
will output
```yaml
a:
  field:
    g: wizz
  fieldA: cat
  fieldB: dog
b:
  field:
    g: wizz
  fieldB: dog
```

## Merge keeps style of LHS
Given a sample.yml file of:
```yaml
a: {things: great}
b:
  also: "me"

```
then
```bash
yq eval '. * {"a":.b}' sample.yml
```
will output
```yaml
a: {things: great, also: "me"}
b:
  also: "me"
```

## Merge arrays
Given a sample.yml file of:
```yaml
a:
  - 1
  - 2
  - 3
b:
  - 3
  - 4
  - 5
```
then
```bash
yq eval '. * {"a":.b}' sample.yml
```
will output
```yaml
a:
  - 3
  - 4
  - 5
b:
  - 3
  - 4
  - 5
```

## Merge, appending arrays
Given a sample.yml file of:
```yaml
a:
  array:
    - 1
    - 2
    - animal: dog
  value: coconut
b:
  array:
    - 3
    - 4
    - animal: cat
  value: banana
```
then
```bash
yq eval '.a *+ .b' sample.yml
```
will output
```yaml
array:
  - 1
  - 2
  - animal: dog
  - 3
  - 4
  - animal: cat
value: banana
```

## Merge to prefix an element
Given a sample.yml file of:
```yaml
a: cat
b: dog
```
then
```bash
yq eval '. * {"a": {"c": .a}}' sample.yml
```
will output
```yaml
a:
  c: cat
b: dog
```

## Merge with simple aliases
Given a sample.yml file of:
```yaml
a: &cat
  c: frog
b:
  f: *cat
c:
  g: thongs
```
then
```bash
yq eval '.c * .b' sample.yml
```
will output
```yaml
g: thongs
f: *cat
```

## Merge does not copy anchor names
Given a sample.yml file of:
```yaml
a:
  c: &cat frog
b:
  f: *cat
c:
  g: thongs
```
then
```bash
yq eval '.c * .a' sample.yml
```
will output
```yaml
g: thongs
c: frog
```

## Merge with merge anchors
Given a sample.yml file of:
```yaml
foo: &foo
  a: foo_a
  thing: foo_thing
  c: foo_c
bar: &bar
  b: bar_b
  thing: bar_thing
  c: bar_c
foobarList:
  b: foobarList_b
  !!merge <<:
    - *foo
    - *bar
  c: foobarList_c
foobar:
  c: foobar_c
  !!merge <<: *foo
  thing: foobar_thing
```
then
```bash
yq eval '.foobar * .foobarList' sample.yml
```
will output
```yaml
c: foobarList_c
<<:
  - *foo
  - *bar
thing: foobar_thing
b: foobarList_b
```

