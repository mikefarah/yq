# Anchor and Alias Operators

Use the `alias` and `anchor` operators to read and write yaml aliases and anchors. The `explode` operator normalises a yaml file (dereference (or expands) aliases and remove anchor names).

`yq` supports merge aliases (like `<<: *blah`) however this is no longer in the standard yaml spec (1.2) and so `yq` will automatically add the `!!merge` tag to these nodes as it is effectively a custom tag.


## NOTE --yaml-fix-merge-anchor-to-spec flag
`yq` doesn't merge anchors `<<:` to spec, in some circumstances it incorrectly overrides existing keys when the spec documents not to do that.

To minimise disruption while still fixing the issue, a flag has been added to toggle this behaviour. This will first default to false; and log warnings to users. Then it will default to true (and still allow users to specify false if needed).

This flag also enables advanced merging, like inline maps, as well as fixes to ensure when exploding a particular path, neighbours are not affect ed.

Long story short, you should be setting this flag to true.

See examples of the flag differences below, where LEGACY is with the flag off; and FIXED is with the flag on.


## Merge one map
see https://yaml.org/type/merge.html

Given a sample.yml file of:
```yaml
- &CENTRE
  x: 1
  y: 2
- &LEFT
  x: 0
  y: 2
- &BIG
  r: 10
- &SMALL
  r: 1
- !!merge <<: *CENTRE
  r: 10
```
then
```bash
yq '.[4] | explode(.)' sample.yml
```
will output
```yaml
x: 1
y: 2
r: 10
```

## Get anchor
Given a sample.yml file of:
```yaml
a: &billyBob cat
```
then
```bash
yq '.a | anchor' sample.yml
```
will output
```yaml
billyBob
```

## Set anchor
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq '.a anchor = "foobar"' sample.yml
```
will output
```yaml
a: &foobar cat
```

## Set anchor relatively using assign-update
Given a sample.yml file of:
```yaml
a:
  b: cat
```
then
```bash
yq '.a anchor |= .b' sample.yml
```
will output
```yaml
a: &cat
  b: cat
```

## Get alias
Given a sample.yml file of:
```yaml
b: &billyBob meow
a: *billyBob
```
then
```bash
yq '.a | alias' sample.yml
```
will output
```yaml
billyBob
```

## Set alias
Given a sample.yml file of:
```yaml
b: &meow purr
a: cat
```
then
```bash
yq '.a alias = "meow"' sample.yml
```
will output
```yaml
b: &meow purr
a: *meow
```

## Set alias to blank does nothing
Given a sample.yml file of:
```yaml
b: &meow purr
a: cat
```
then
```bash
yq '.a alias = ""' sample.yml
```
will output
```yaml
b: &meow purr
a: cat
```

## Set alias relatively using assign-update
Given a sample.yml file of:
```yaml
b: &meow purr
a:
  f: meow
```
then
```bash
yq '.a alias |= .f' sample.yml
```
will output
```yaml
b: &meow purr
a: *meow
```

## Explode alias and anchor
Given a sample.yml file of:
```yaml
f:
  a: &a cat
  b: *a
```
then
```bash
yq 'explode(.f)' sample.yml
```
will output
```yaml
f:
  a: cat
  b: cat
```

## Explode with no aliases or anchors
Given a sample.yml file of:
```yaml
a: mike
```
then
```bash
yq 'explode(.a)' sample.yml
```
will output
```yaml
a: mike
```

## Explode with alias keys
Given a sample.yml file of:
```yaml
f:
  a: &a cat
  *a : b
```
then
```bash
yq 'explode(.f)' sample.yml
```
will output
```yaml
f:
  a: cat
  cat: b
```

## Dereference and update a field
Use explode with multiply to dereference an object

Given a sample.yml file of:
```yaml
item_value: &item_value
  value: true
thingOne:
  name: item_1
  !!merge <<: *item_value
thingTwo:
  name: item_2
  !!merge <<: *item_value
```
then
```bash
yq '.thingOne |= (explode(.) | sort_keys(.)) * {"value": false}' sample.yml
```
will output
```yaml
item_value: &item_value
  value: true
thingOne:
  name: item_1
  value: false
thingTwo:
  name: item_2
  !!merge <<: *item_value
```

## LEGACY: Explode with merge anchors
Caution: this is for when --yaml-fix-merge-anchor-to-spec=false; it's not to YAML spec because the merge anchors incorrectly override the object values (foobarList.b is set to bar_b when it should still be foobarList_b). Flag will default to true in late 2025

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
yq 'explode(.)' sample.yml
```
will output
```yaml
foo:
  a: foo_a
  thing: foo_thing
  c: foo_c
bar:
  b: bar_b
  thing: bar_thing
  c: bar_c
foobarList:
  b: bar_b
  thing: foo_thing
  c: foobarList_c
  a: foo_a
foobar:
  c: foo_c
  a: foo_a
  thing: foobar_thing
```

## LEGACY: Merge multiple maps
see https://yaml.org/type/merge.html. This has the correct data, but the wrong key order; set --yaml-fix-merge-anchor-to-spec=true to fix the key order.

Given a sample.yml file of:
```yaml
- &CENTRE
  x: 1
  y: 2
- &LEFT
  x: 0
  y: 2
- &BIG
  r: 10
- &SMALL
  r: 1
- !!merge <<:
    - *CENTRE
    - *BIG
```
then
```bash
yq '.[4] | explode(.)' sample.yml
```
will output
```yaml
r: 10
x: 1
y: 2
```

## LEGACY: Override
see https://yaml.org/type/merge.html. This has the correct data, but the wrong key order; set --yaml-fix-merge-anchor-to-spec=true to fix the key order.

Given a sample.yml file of:
```yaml
- &CENTRE
  x: 1
  y: 2
- &LEFT
  x: 0
  y: 2
- &BIG
  r: 10
- &SMALL
  r: 1
- !!merge <<:
    - *BIG
    - *LEFT
    - *SMALL
  x: 1
```
then
```bash
yq '.[4] | explode(.)' sample.yml
```
will output
```yaml
r: 10
x: 1
y: 2
```

## FIXED: Explode with merge anchors
Set `--yaml-fix-merge-anchor-to-spec=true` to get this correct merge behaviour (flag will default to true in late 2025).
Observe that foobarList.b property is still foobarList_b.

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
yq 'explode(.)' sample.yml
```
will output
```yaml
foo:
  a: foo_a
  thing: foo_thing
  c: foo_c
bar:
  b: bar_b
  thing: bar_thing
  c: bar_c
foobarList:
  b: foobarList_b
  a: foo_a
  thing: foo_thing
  c: foobarList_c
foobar:
  c: foobar_c
  a: foo_a
  thing: foobar_thing
```

## FIXED: Merge multiple maps
Set `--yaml-fix-merge-anchor-to-spec=true` to get this correct merge behaviour (flag will default to true in late 2025).
Taken from https://yaml.org/type/merge.html. Same values as legacy, but with the correct key order.

Given a sample.yml file of:
```yaml
- &CENTRE
  x: 1
  y: 2
- &LEFT
  x: 0
  y: 2
- &BIG
  r: 10
- &SMALL
  r: 1
- !!merge <<:
    - *CENTRE
    - *BIG
```
then
```bash
yq '.[4] | explode(.)' sample.yml
```
will output
```yaml
x: 1
y: 2
r: 10
```

## FIXED: Override
Set `--yaml-fix-merge-anchor-to-spec=true` to get this correct merge behaviour (flag will default to true in late 2025).
Taken from https://yaml.org/type/merge.html. Same values as legacy, but with the correct key order.

Given a sample.yml file of:
```yaml
- &CENTRE
  x: 1
  y: 2
- &LEFT
  x: 0
  y: 2
- &BIG
  r: 10
- &SMALL
  r: 1
- !!merge <<:
    - *BIG
    - *LEFT
    - *SMALL
  x: 1
```
then
```bash
yq '.[4] | explode(.)' sample.yml
```
will output
```yaml
r: 10
y: 2
x: 1
```

## Exploding inline merge anchor
Set `--yaml-fix-merge-anchor-to-spec=true` to get this correct merge behaviour (flag will default to true in late 2025).


Given a sample.yml file of:
```yaml
a:
  b: &b 42
!!merge <<:
  c: *b
```
then
```bash
yq 'explode(.) | sort_keys(.)' sample.yml
```
will output
```yaml
a:
  b: 42
c: 42
```

