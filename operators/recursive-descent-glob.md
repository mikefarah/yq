# Recursive Descent (Glob)

This operator recursively matches (or globs) all children nodes given of a particular element, including that node itself. This is most often used to apply a filter recursively against all matches. It can be used in either the

## match values form `..`

This will, like the `jq` equivalent, recursively match all _value_ nodes. Use it to find/manipulate particular values.

For instance to set the `style` of all _value_ nodes in a yaml doc, excluding map keys:

```bash
yq eval '.. style= "flow"' file.yaml
```

## match values and map keys form `...`

The also includes map keys in the results set. This is particularly useful in YAML as unlike JSON, map keys can have their own styling, tags and use anchors and aliases.

For instance to set the `style` of all nodes in a yaml doc, including the map keys:

```bash
yq eval '... style= "flow"' file.yaml
```

## Recurse map (values only)

Given a sample.yml file of:

```yaml
a: frog
```

then

```bash
yq eval '..' sample.yml
```

will output

```yaml
a: frog
frog
```

## Recursively find nodes with keys

Note that this example has wrapped the expression in `[]` to show that there are two matches returned. You do not have to wrap in `[]` in your path expression.

Given a sample.yml file of:

```yaml
a:
  name: frog
  b:
    name: blog
    age: 12
```

then

```bash
yq eval '[.. | select(has("name"))]' sample.yml
```

will output

```yaml
- name: frog
  b:
    name: blog
    age: 12
- name: blog
  age: 12
```

## Recursively find nodes with values

Given a sample.yml file of:

```yaml
a:
  nameA: frog
  b:
    nameB: frog
    age: 12
```

then

```bash
yq eval '.. | select(. == "frog")' sample.yml
```

will output

```yaml
frog
frog
```

## Recurse map (values and keys)

Note that the map key appears in the results

Given a sample.yml file of:

```yaml
a: frog
```

then

```bash
yq eval '...' sample.yml
```

will output

```yaml
a: frog
a
frog
```

## Aliases are not traversed

Given a sample.yml file of:

```yaml
a: &cat
  c: frog
b: *cat
```

then

```bash
yq eval '[..]' sample.yml
```

will output

```yaml
- a: &cat
    c: frog
  b: *cat
- &cat
  c: frog
- frog
- *cat
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
yq eval '.foobar | [..]' sample.yml
```

will output

```yaml
- c: foobar_c
  !!merge <<: *foo
  thing: foobar_thing
- foobar_c
- *foo
- foobar_thing
```
