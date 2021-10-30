# Encoder / Decoder

Encode operators will take the piped in object structure and encode it as a string in the desired format. The decode operators do the opposite, they take a formatted string and decode it into the relevant object structure.

Note that you can optionally pass an indent value to the encode functions (see below).

These operators are useful to process yaml documents that have stringified embeded yaml/json/props in them.

## Encode value as yaml string
Indent defaults to 2

Given a sample.yml file of:
```yaml
a:
  cool:
    bob: dylan
```
then
```bash
yq eval '.b = (.a | to_yaml)' sample.yml
```
will output
```yaml
a:
  cool:
    bob: dylan
b: |
  cool:
    bob: dylan
```

## Encode value as yaml string, with custom indentation
You can specify the indentation level as the first parameter.

Given a sample.yml file of:
```yaml
a:
  cool:
    bob: dylan
```
then
```bash
yq eval '.b = (.a | to_yaml(8))' sample.yml
```
will output
```yaml
a:
  cool:
    bob: dylan
b: |
  cool:
          bob: dylan
```

## Encode value as yaml string, using toyaml
Does the same thing as to_yaml, matching jq naming convention.

Given a sample.yml file of:
```yaml
a:
  cool: thing
```
then
```bash
yq eval '.b = (.a | to_yaml)' sample.yml
```
will output
```yaml
a:
  cool: thing
b: |
  cool: thing
```

## Encode value as json string
Given a sample.yml file of:
```yaml
a:
  cool: thing
```
then
```bash
yq eval '.b = (.a | to_json)' sample.yml
```
will output
```yaml
a:
  cool: thing
b: |
  {
    "cool": "thing"
  }
```

## Encode value as json string, on one line
Pass in a 0 indent to print json on a single line.

Given a sample.yml file of:
```yaml
a:
  cool: thing
```
then
```bash
yq eval '.b = (.a | to_json(0))' sample.yml
```
will output
```yaml
a:
  cool: thing
b: '{"cool":"thing"}'
```

## Encode value as props string
Given a sample.yml file of:
```yaml
a:
  cool: thing
```
then
```bash
yq eval '.b = (.a | to_props)' sample.yml
```
will output
```yaml
a:
  cool: thing
b: |
  cool = thing
```

## Decode a yaml encoded string
Given a sample.yml file of:
```yaml
a: 'foo: bar'
```
then
```bash
yq eval '.b = (.a | from_yaml)' sample.yml
```
will output
```yaml
a: 'foo: bar'
b:
  foo: bar
```

## Update a multiline encoded yaml string
Given a sample.yml file of:
```yaml
a: |
  foo: bar
  baz: dog

```
then
```bash
yq eval '.a |= (from_yaml | .foo = "cat" | to_yaml)' sample.yml
```
will output
```yaml
a: |
  foo: cat
  baz: dog
```

## Update a single line encoded yaml string
Given a sample.yml file of:
```yaml
a: 'foo: bar'
```
then
```bash
yq eval '.a |= (from_yaml | .foo = "cat" | to_yaml)' sample.yml
```
will output
```yaml
a: 'foo: cat'
```

