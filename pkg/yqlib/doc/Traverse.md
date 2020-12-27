This is the simplest (and perhaps most used) operator, it is used to navigate deeply into yaml structurse.
## Simple map navigation
Given a sample.yml file of:
```yaml
a:
  b: apple
```
then
```bash
yq eval '.a' sample.yml
```
will output
```yaml
b: apple
```

## Splat
Often used to pipe children into other operators

Given a sample.yml file of:
```yaml
- b: apple
- c: banana
```
then
```bash
yq eval '.[]' sample.yml
```
will output
```yaml
b: apple
c: banana
```

## Special characters
Use quotes around path elements with special characters

Given a sample.yml file of:
```yaml
"{}": frog
```
then
```bash
yq eval '."{}"' sample.yml
```
will output
```yaml
frog
```

## Children don't exist
Nodes are added dynamically while traversing

Given a sample.yml file of:
```yaml
c: banana
```
then
```bash
yq eval '.a.b' sample.yml
```
will output
```yaml
null
```

## Wildcard matching
Given a sample.yml file of:
```yaml
a:
  cat: apple
  mad: things
```
then
```bash
yq eval '.a."*a*"' sample.yml
```
will output
```yaml
apple
things
```

## Aliases
Given a sample.yml file of:
```yaml
a: &cat
  c: frog
b: *cat
```
then
```bash
yq eval '.b' sample.yml
```
will output
```yaml
*cat
```

## Traversing aliases with splat
Given a sample.yml file of:
```yaml
a: &cat
  c: frog
b: *cat
```
then
```bash
yq eval '.b.[]' sample.yml
```
will output
```yaml
frog
```

## Traversing aliases explicitly
Given a sample.yml file of:
```yaml
a: &cat
  c: frog
b: *cat
```
then
```bash
yq eval '.b.c' sample.yml
```
will output
```yaml
frog
```

## Traversing arrays by index
Given a sample.yml file of:
```yaml
- 1
- 2
- 3
```
then
```bash
yq eval '.[0]' sample.yml
```
will output
```yaml
1
```

## Maps with numeric keys
Given a sample.yml file of:
```yaml
2: cat
```
then
```bash
yq eval '.[2]' sample.yml
```
will output
```yaml
cat
```

## Maps with non existing numeric keys
Given a sample.yml file of:
```yaml
a: b
```
then
```bash
yq eval '.[0]' sample.yml
```
will output
```yaml
null
```

## Traversing merge anchors
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
yq eval '.foobar.a' sample.yml
```
will output
```yaml
foo_a
```

## Traversing merge anchors with override
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
yq eval '.foobar.c' sample.yml
```
will output
```yaml
foo_c
```

## Traversing merge anchors with local override
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
yq eval '.foobar.thing' sample.yml
```
will output
```yaml
foobar_thing
```

## Splatting merge anchors
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
yq eval '.foobar.[]' sample.yml
```
will output
```yaml
foobar_c
*foo
foobar_thing
```

## Traversing merge anchor lists
Note that the later merge anchors override previous

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
yq eval '.foobarList.thing' sample.yml
```
will output
```yaml
bar_thing
```

## Splatting merge anchor lists
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
yq eval '.foobarList.[]' sample.yml
```
will output
```yaml
foobarList_b
- *foo
- *bar
foobarList_c
```

## Select multiple indices
Given a sample.yml file of:
```yaml
a:
  - a
  - b
  - c
```
then
```bash
yq eval '.a[0, 2]' sample.yml
```
will output
```yaml
a
c
```

