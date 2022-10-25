# Multiply (Merge)

Like the multiple operator in jq, depending on the operands, this multiply operator will do different things. Currently numbers, arrays and objects are supported.

## Objects and arrays - merging
Objects are merged deeply matching on matching keys. By default, array values override and are not deeply merged.

Note that when merging objects, this operator returns the merged object (not the parent). This will be clearer in the examples below.

### Merge Flags
You can control how objects are merged by using one or more of the following flags. Multiple flags can be used together, e.g. `.a *+? .b`.  See examples below

- `+` append arrays
- `d` deeply merge arrays
- `?` only merge _existing_ fields
- `n` only merge _new_ fields
- `c` clobber custom tags


### Merge two files together
This uses the load operator to merge file2 into file1.
```bash
yq '. *= load("file2.yml")' file1.yml
```

### Merging all files
Note the use of `eval-all` to ensure all documents are loaded into memory.

```bash
yq eval-all '. as $item ireduce ({}; . * $item )' *.yml
```

# Merging complex arrays together by a key field
By default - `yq` merge is naive. It merges maps when they match the key name, and arrays are merged either by appending them together, or merging the entries by their position in the array.

For more complex array merging (e.g. merging items that match on a certain key) please see the example [here](https://mikefarah.gitbook.io/yq/operators/multiply-merge#merge-arrays-of-objects-together-matching-on-a-key)


## Multiply integers
Given a sample.yml file of:
```yaml
a: 3
b: 4
```
then
```bash
yq '.a *= .b' sample.yml
```
will output
```yaml
a: 12
b: 4
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
yq '.a * .b' sample.yml
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
yq '. * {"a":.b}' sample.yml
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
yq '. * {"a":.b}' sample.yml
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
yq '. * {"a":.b}' sample.yml
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
yq '.a *? .b' sample.yml
```
will output
```yaml
thing: two
cat: frog
```

## Merge, only new fields
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
yq '.a *n .b' sample.yml
```
will output
```yaml
thing: one
cat: frog
missing: two
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
yq '.a *+ .b' sample.yml
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
yq '.a *?+ .b' sample.yml
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
yq '.a *d .b' sample.yml
```
will output
```yaml
- name: fred
  age: 34
- name: bob
  age: 32
```

## Merge arrays of objects together, matching on a key

This is a fairly complex expression - you can use it as is by providing the environment variables as seen in the example below.

It merges in the array provided in the second file into the first - matching on equal keys.

Explanation:

The approach, at a high level, is to reduce into a merged map (keyed by the unique key)
and then convert that back into an array.

First the expression will create a map from the arrays keyed by the idPath, the unique field we want to merge by.
The reduce operator is merging '({}; . * $item )', so array elements with the matching key will be merged together.

Next, we convert the map back to an array, using reduce again, concatenating all the map values together.

Finally, we set the result of the merged array back into the first doc.

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
newArray:
  - a: banana
    c: bananaC
  - a: apple
    b: appleB2
  - a: dingo
    c: dingoC
```
then
```bash
idPath=".a"  originalPath=".myArray"  otherPath=".newArray" yq eval-all '
(
  (( (eval(strenv(originalPath)) + eval(strenv(otherPath)))  | .[] | {(eval(strenv(idPath))):  .}) as $item ireduce ({}; . * $item )) as $uniqueMap
  | ( $uniqueMap  | to_entries | .[]) as $item ireduce([]; . + $item.value)
) as $mergedArray
| select(fi == 0) | (eval(strenv(originalPath))) = $mergedArray
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
yq '. * {"a": {"c": .a}}' sample.yml
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
yq '.c * .b' sample.yml
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
yq '.c * .a' sample.yml
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
yq '.foobar * .foobarList' sample.yml
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

## Custom types: that are really numbers
When custom tags are encountered, yq will try to decode the underlying type.

Given a sample.yml file of:
```yaml
a: !horse 2
b: !goat 3
```
then
```bash
yq '.a = .a * .b' sample.yml
```
will output
```yaml
a: !horse 6
b: !goat 3
```

## Custom types: that are really maps
Custom tags will be maintained.

Given a sample.yml file of:
```yaml
a: !horse
  cat: meow
b: !goat
  dog: woof
```
then
```bash
yq '.a = .a * .b' sample.yml
```
will output
```yaml
a: !horse
  cat: meow
  dog: woof
b: !goat
  dog: woof
```

## Custom types: clobber tags
Use the `c` option to clobber custom tags. Note that the second tag is now used

Given a sample.yml file of:
```yaml
a: !horse
  cat: meow
b: !goat
  dog: woof
```
then
```bash
yq '.a *=c .b' sample.yml
```
will output
```yaml
a: !goat
  cat: meow
  dog: woof
b: !goat
  dog: woof
```

