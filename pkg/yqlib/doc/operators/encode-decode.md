# Encoder / Decoder

Encode operators will take the piped in object structure and encode it as a string in the desired format. The decode operators do the opposite, they take a formatted string and decode it into the relevant object structure.

Note that you can optionally pass an indent value to the encode functions (see below).

These operators are useful to process yaml documents that have stringified embedded yaml/json/props in them.


| Format | Decode (from string) | Encode (to string) |
| --- | -- | --|
| Yaml | from_yaml/@yamld | to_yaml(i)/@yaml |
| JSON | from_json/@jsond | to_json(i)/@json |
| Properties | from_props/@propsd  | to_props/@props |
| CSV | from_csv/@csvd | to_csv/@csv |
| TSV | from_tsv/@tsvd | to_tsv/@tsv |
| XML | from_xml/@xmld | to_xml(i)/@xml |
| Base64 | @base64d | @base64 |
| URI | @urid | @uri |
| Shell |  | @sh |


See CSV and TSV [documentation](https://mikefarah.gitbook.io/yq/usage/csv-tsv) for accepted formats.

XML uses the `--xml-attribute-prefix` and `xml-content-name` flags to identify attributes and content fields.


Base64 assumes [rfc4648](https://rfc-editor.org/rfc/rfc4648.html) encoding. Encoding and decoding both assume that the content is a utf-8 string and not binary content.

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

## Decode props encoded string
Given a sample.yml file of:
```yaml
a: |-
  cats=great
  dogs=cool as well
```
then
```bash
yq '.a |= @propsd' sample.yml
```
will output
```yaml
a:
  cats: great
  dogs: cool as well
```

## Decode csv encoded string
Given a sample.yml file of:
```yaml
a: |-
  cats,dogs
  great,cool as well
```
then
```bash
yq '.a |= @csvd' sample.yml
```
will output
```yaml
a:
  - cats: great
    dogs: cool as well
```

## Decode tsv encoded string
Given a sample.yml file of:
```yaml
a: |-
  cats	dogs
  great	cool as well
```
then
```bash
yq '.a |= @tsvd' sample.yml
```
will output
```yaml
a:
  - cats: great
    dogs: cool as well
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

## Encode array of arrays as tsv string
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
    +@id: hi
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
    +@id: hi
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
    +@id: hi
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

## Encode a string to base64
Given a sample.yml file of:
```yaml
coolData: a special string
```
then
```bash
yq '.coolData | @base64' sample.yml
```
will output
```yaml
YSBzcGVjaWFsIHN0cmluZw==
```

## Encode a yaml document to base64
Pipe through @yaml first to convert to a string, then use @base64 to encode it.

Given a sample.yml file of:
```yaml
a: apple
```
then
```bash
yq '@yaml | @base64' sample.yml
```
will output
```yaml
YTogYXBwbGUK
```

## Encode a string to uri
Given a sample.yml file of:
```yaml
coolData: this has & special () characters *
```
then
```bash
yq '.coolData | @uri' sample.yml
```
will output
```yaml
this+has+%26+special+%28%29+characters+%2A
```

## Decode a URI to a string
Given a sample.yml file of:
```yaml
this+has+%26+special+%28%29+characters+%2A
```
then
```bash
yq '@urid' sample.yml
```
will output
```yaml
this has & special () characters *
```

## Encode a string to sh
Sh/Bash friendly string

Given a sample.yml file of:
```yaml
coolData: strings with spaces and a 'quote'
```
then
```bash
yq '.coolData | @sh' sample.yml
```
will output
```yaml
strings' with spaces and a '\'quote\'
```

## Decode a base64 encoded string
Decoded data is assumed to be a string.

Given a sample.yml file of:
```yaml
coolData: V29ya3Mgd2l0aCBVVEYtMTYg8J+Yig==
```
then
```bash
yq '.coolData | @base64d' sample.yml
```
will output
```yaml
Works with UTF-16 😊
```

## Decode a base64 encoded yaml document
Pipe through `from_yaml` to parse the decoded base64 string as a yaml document.

Given a sample.yml file of:
```yaml
coolData: YTogYXBwbGUK
```
then
```bash
yq '.coolData |= (@base64d | from_yaml)' sample.yml
```
will output
```yaml
coolData:
  a: apple
```

