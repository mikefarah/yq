This operator recursively matches all children nodes given of a particular element, including that node itself. This is most often used to apply a filter recursively against all matches, for instance to set the `style` of all nodes in a yaml doc:

```bash
yq eval '.. style= "flow"' file.yaml
```
## Aliases are not traversed
Given a sample.yml file of:
```yaml
a: &cat {c: frog}
b: *cat
'': null
```
then
```bash
yq eval '[..]' sample.yml
```
will output
```yaml
- a: &cat {c: frog}
  b: *cat
  '': null
- null
```

## Merge docs are not traversed
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
yq eval '.foobar | [..]' sample.yml
```
will output
```yaml
- c: foobar_c
  !!merge <<: *foo
  thing: foobar_thing
  '': null
- null
```

