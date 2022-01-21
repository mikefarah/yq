# JSON

Encode and decode to and from JSON. Note that YAML is a _superset_ of JSON - so `yq` can read any json file without doing anything special.

This means you don't need to 'convert' a JSON file to YAML - however if you want idiomatic YAML styling, then you can use the `-P/--prettyPrint` flag, see examples below.

## Parse json: simple
JSON is a subset of yaml, so all you need to do is prettify the output

Given a sample.json file of:
```json
{"cat": "meow"}
```
then
```bash
yq e -P '.' sample.json
```
will output
```yaml
cat: meow
```

## Parse json: complex
JSON is a subset of yaml, so all you need to do is prettify the output

Given a sample.json file of:
```json
{"a":"Easy! as one two three","b":{"c":2,"d":[3,4]}}
```
then
```bash
yq e -P '.' sample.json
```
will output
```yaml
a: Easy! as one two three
b:
  c: 2
  d:
    - 3
    - 4
```

## Encode json: simple
Given a sample.yml file of:
```yaml
cat: meow
```
then
```bash
yq e -o=json '.' sample.yml
```
will output
```json
{
  "cat": "meow"
}
```

## Encode json: simple - in one line
Given a sample.yml file of:
```yaml
cat: meow # this is a comment, and it will be dropped.
```
then
```bash
yq e -o=json -I=0 '.' sample.yml
```
will output
```json
{"cat":"meow"}
```

## Encode json: comments
Given a sample.yml file of:
```yaml
cat: meow # this is a comment, and it will be dropped.
```
then
```bash
yq e -o=json '.' sample.yml
```
will output
```json
{
  "cat": "meow"
}
```

## Encode json: anchors
Anchors are dereferenced

Given a sample.yml file of:
```yaml
cat: &ref meow
anotherCat: *ref
```
then
```bash
yq e -o=json '.' sample.yml
```
will output
```json
{
  "cat": "meow",
  "anotherCat": "meow"
}
```

## Encode json: multiple results
Each matching node is converted into a json doc. This is best used with 0 indent (json document per line)

Given a sample.yml file of:
```yaml
things: [{stuff: cool}, {whatever: cat}]
```
then
```bash
yq e -o=json -I=0 '.things[]' sample.yml
```
will output
```json
{"stuff":"cool"}
{"whatever":"cat"}
```

