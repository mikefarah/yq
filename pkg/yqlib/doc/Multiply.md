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
a: {field: me, fieldA: cat}
b: {field: {g: wizz}, fieldB: dog}
'': null
```
then
```bash
yq eval '.a * .b' sample.yml
```
will output
```yaml
{'': null}
```

## Merge objects together, returning parent object
Given a sample.yml file of:
```yaml
a: {field: me, fieldA: cat}
b: {field: {g: wizz}, fieldB: dog}
'': null
```
then
```bash
yq eval '. * {"a":.b}' sample.yml
```
will output
```yaml
'': null
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
'': null
```

## Merge arrays
Given a sample.yml file of:
```yaml
a: [1, 2, 3]
b: [3, 4, 5]
'': null
```
then
```bash
yq eval '. * {"a":.b}' sample.yml
```
will output
```yaml
'': null
```

## Merge, appending arrays
Given a sample.yml file of:
```yaml
a: {array: [1, 2, {animal: dog}], value: coconut}
b: {array: [3, 4, {animal: cat}], value: banana}
'': null
```
then
```bash
yq eval '.a *+ .b' sample.yml
```
will output
```yaml
{'': null}
```

## Merge to prefix an element
Given a sample.yml file of:
```yaml
a: cat
b: dog
'': null
```
then
```bash
yq eval '. * {"a": {"c": .a}}' sample.yml
```
will output
```yaml
'': null
```

## Merge with simple aliases
Given a sample.yml file of:
```yaml
a: &cat {c: frog}
b: {f: *cat}
c: {g: thongs}
'': null
```
then
```bash
yq eval '.c * .b' sample.yml
```
will output
```yaml
{'': null}
```

## Merge does not copy anchor names
Given a sample.yml file of:
```yaml
a: {c: &cat frog}
b: {f: *cat}
c: {g: thongs}
'': null
```
then
```bash
yq eval '.c * .a' sample.yml
```
will output
```yaml
{'': null}
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
  !!merge <<: [*foo, *bar]
  c: foobarList_c
foobar:
  c: foobar_c
  !!merge <<: *foo
  thing: foobar_thing
'': null
```
then
```bash
yq eval '.foobar * .foobarList' sample.yml
```
will output
```yaml
'': null
```

