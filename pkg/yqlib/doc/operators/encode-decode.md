# Encoder / Decoder

Encode operators will take the piped in object structure and encode it as a string in the desired format. The decode operators do the opposite, they take a formatted string and decode it into the relevant object structure.

Note that you can optionally pass an indent value to the encode functions (see below).

These operators are useful to process yaml documents that have stringified embeded yaml/json/props in them.


| Format | Decode (from string) | Encode (to string) |
| --- | -- | --|
| Yaml | from_yaml | to_yaml(i)/@yaml |
| JSON | from_json | to_json(i)/@json |
| Properties |  | to_props/@props |
| CSV |  | to_csv/@csv |
| TSV |  | to_tsv/@tsv |
| XML | from_xml | to_xml(i)/@xml |


CSV and TSV format both accept either a single array or scalars (representing a single row), or an array of array of scalars (representing multiple rows). 

XML uses the `--xml-attribute-prefix` and `xml-content-name` flags to identify attributes and content fields.


## Encode value as json string
Given a sample.yml file of:
```yaml
a:
  cool: thing
```
then
```bash
yq '.b = (.a | to_json)' sample.yml
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
yq '.b = (.a | to_json(0))' sample.yml
```
will output
```yaml
a:
  cool: thing
b: '{"cool":"thing"}'
```

## Encode value as json string, on one line shorthand
Pass in a 0 indent to print json on a single line.

Given a sample.yml file of:
```yaml
a:
  cool: thing
```
then
```bash
yq '.b = (.a | @json)' sample.yml
```
will output
```yaml
a:
  cool: thing
b: '{"cool":"thing"}'
```

## Decode a json encoded string
Keep in mind JSON is a subset of YAML. If you want idiomatic yaml, pipe through the style operator to clear out the JSON styling.

Given a sample.yml file of:
```yaml
a: '{"cool":"thing"}'
```
then
```bash
yq '.a | from_json | ... style=""' sample.yml
```
will output
```yaml
cool: thing
```

## Encode value as props string
Given a sample.yml file of:
```yaml
a:
  cool: thing
```
then
```bash
yq '.b = (.a | @props)' sample.yml
```
will output
```yaml
a:
  cool: thing
b: |
  cool = thing
```

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
yq '.b = (.a | to_yaml)' sample.yml
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
yq '.b = (.a | to_yaml(8))' sample.yml
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

## Decode a yaml encoded string
Given a sample.yml file of:
```yaml
a: 'foo: bar'
```
then
```bash
yq '.b = (.a | from_yaml)' sample.yml
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
yq '.a |= (from_yaml | .foo = "cat" | to_yaml)' sample.yml
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
yq '.a |= (from_yaml | .foo = "cat" | to_yaml)' sample.yml
```
will output
```yaml
a: 'foo: cat'
```

## Encode array of scalars as csv string
Scalars are strings, numbers and booleans.

Given a sample.yml file of:
```yaml
- cat
- thing1,thing2
- true
- 3.40
```
then
```bash
yq '@csv' sample.yml
```
will output
```yaml
cat,"thing1,thing2",true,3.40
```

## Encode array of arrays as csv string
Given a sample.yml file of:
```yaml
- - cat
  - thing1,thing2
  - true
  - 3.40
- - dog
  - thing3
  - false
  - 12
```
then
```bash
yq '@csv' sample.yml
```
will output
```yaml
cat,"thing1,thing2",true,3.40
dog,thing3,false,12
```

## Encode array of array scalars as tsv string
Scalars are strings, numbers and booleans.

Given a sample.yml file of:
```yaml
- - cat
  - thing1,thing2
  - true
  - 3.40
- - dog
  - thing3
  - false
  - 12
```
then
```bash
yq '@tsv' sample.yml
```
will output
```yaml
cat	thing1,thing2	true	3.40
dog	thing3	false	12
```

## Encode value as xml string
Given a sample.yml file of:
```yaml
a:
  cool:
    foo: bar
    +id: hi
```
then
```bash
yq '.a | to_xml' sample.yml
```
will output
```yaml
<cool id="hi">
  <foo>bar</foo>
</cool>

```

## Encode value as xml string on a single line
Given a sample.yml file of:
```yaml
a:
  cool:
    foo: bar
    +id: hi
```
then
```bash
yq '.a | @xml' sample.yml
```
will output
```yaml
<cool id="hi"><foo>bar</foo></cool>

```

## Encode value as xml string with custom indentation
Given a sample.yml file of:
```yaml
a:
  cool:
    foo: bar
    +id: hi
```
then
```bash
yq '{"cat": .a | to_xml(1)}' sample.yml
```
will output
```yaml
cat: |
  <cool id="hi">
   <foo>bar</foo>
  </cool>
```

## Decode a xml encoded string
Given a sample.yml file of:
```yaml
a: <foo>bar</foo>
```
then
```bash
yq '.b = (.a | from_xml)' sample.yml
```
will output
```yaml
a: <foo>bar</foo>
b:
  foo: bar
```

