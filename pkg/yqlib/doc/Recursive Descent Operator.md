This operator recursively matches all children nodes given of a particular element, including that node itself. This is most often used to apply a filter recursively against all matches, for instance to set the `style` of all nodes in a yaml doc:

```bash
yq eval '.. style= "flow"' file.yaml
```
## Examples
### Map
Given a sample.yml file of:
```yaml
a:
  b: apple
```
then
```bash
yq eval '..' sample.yml
```
will output
```yaml
a:
  b: apple
b: apple
apple
```

### Array
Given a sample.yml file of:
```yaml
- 1
- 2
- 3
```
then
```bash
yq eval '..' sample.yml
```
will output
```yaml
- 1
- 2
- 3
1
2
3
```

### Array of maps
Given a sample.yml file of:
```yaml
- a: cat
- 2
- true
```
then
```bash
yq eval '..' sample.yml
```
will output
```yaml
- a: cat
- 2
- true
a: cat
cat
2
true
```

### Aliases are not traversed
Given a sample.yml file of:
```yaml
a: &cat
  c: frog
b: *cat
```
then
```bash
yq eval '..' sample.yml
```
will output
```yaml
a: &cat
  c: frog
b: *cat
&cat
c: frog
frog
*cat
```

### Merge docs are not traversed
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
yq eval '.foobar | ..' sample.yml
```
will output
```yaml
c: foobar_c
!!merge <<: *foo
thing: foobar_thing
foobar_c
*foo
foobar_thing
```

