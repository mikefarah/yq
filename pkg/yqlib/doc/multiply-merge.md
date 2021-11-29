# Multiply (Merge)

Like the multiple operator in jq, depending on the operands, this multiply operator will do different things. Currently numbers, arrays and objects are supported.

## Objects and arrays - merging
Objects are merged deeply matching on matching keys. By default, array values override and are not deeply merged.

Note that when merging objects, this operator returns the merged object (not the parent). This will be clearer in the examples below.

### Merge Flags
You can control how objects are merged by using one or more of the following flags. Multiple flags can be used together, e.g. `.a *+? .b`.  See examples below

- `+` to append arrays
- `?` to only merge existing fields
- `d` to deeply merge arrays

### Merging files
Note the use of `eval-all` to ensure all documents are loaded into memory.

```bash
yq eval-all 'select(fileIndex == 0) * select(fileIndex == 1)' file1.yaml file2.yaml
```

## Multiply integers
Running
```bash
yq eval --null-input '3 * 4'
```
will output
```yaml
12
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

## Merge, only existing fields
Given a sample.yml file of:
```yaml
a:
  thing: one
  cat: frog
b:
  missing: two
  thing: two
```
then
```bash
yq eval '.a *? .b' sample.yml
```
will output
```yaml
thing: two
cat: frog
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

## Merge, only existing fields, appending arrays
Given a sample.yml file of:
```yaml
a:
  thing:
    - 1
    - 2
b:
  thing:
    - 3
    - 4
  another:
    - 1
```
then
```bash
yq eval '.a *?+ .b' sample.yml
```
will output
```yaml
thing:
  - 1
  - 2
  - 3
  - 4
```

## Merge, deeply merging arrays
Merging arrays deeply means arrays are merge like objects, with indexes as their key. In this case, we merge the first item in the array, and do nothing with the second.

Given a sample.yml file of:
```yaml
a:
  - name: fred
    age: 12
  - name: bob
    age: 32
b:
  - name: fred
    age: 34
```
then
```bash
yq eval '.a *d .b' sample.yml
```
will output
```yaml
- name: fred
  age: 34
- name: bob
  age: 32
```

## Merge arrays of objects together, matching on a key
There are several parts of the complex expression. 
The first part is doing the hard work, it creates a map from the arrays keyed by '.a', so that there are no duplicates. 
Then there's another reduce that converts that map back to an array.
Finally, we set the result of the merged array back into the first doc.

To use this, you will need to update '.myArray' to be the expression to your array (e.g. .my.array), and '.a' to be the key field of your array (e.g. '.name')

Thanks Kev from [stackoverflow](https://stackoverflow.com/a/70109529/1168223)


Given a sample.yml file of:
```yaml
myArray:
  - a: apple
    b: appleB
  - a: kiwi
    b: kiwiB
  - a: banana
    b: bananaB
something: else
```
And another sample another.yml file of:
```yaml
myArray:
  - a: banana
    c: bananaC
  - a: apple
    b: appleB2
  - a: dingo
    c: dingoC
```
then
```bash
yq eval-all '
(
  ((.myArray[] | {.a: .}) as $item ireduce ({}; . * $item )) as $uniqueMap
  | ( $uniqueMap  | to_entries | .[]) as $item ireduce([]; . + $item.value)
) as $mergedArray
| select(fi == 0) | .myArray = $mergedArray
' sample.yml another.yml
```
will output
```yaml
myArray:
  - a: apple
    b: appleB2
  - a: kiwi
    b: kiwiB
  - a: banana
    b: bananaB
    c: bananaC
  - a: dingo
    c: dingoC
something: else
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

## Merge copies anchor names
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
c: &cat frog
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
!!merge <<:
  - *foo
  - *bar
thing: foobar_thing
b: foobarList_b
```

