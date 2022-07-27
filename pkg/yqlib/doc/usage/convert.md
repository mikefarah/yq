# JSON

Encode and decode to and from JSON. Note that, unless you have multiple JSON documents in a single file, YAML is a _superset_ of JSON - so `yq` can read any json file without doing anything special.

If you do have mulitple JSON documents in a single file (e.g. NDJSON) then you will need to explicity use the json parser `-p=json`.

This means you don't need to 'convert' a JSON file to YAML - however if you want idiomatic YAML styling, then you can use the `-P/--prettyPrint` flag, see examples below.

{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Parse json: simple
JSON is a subset of yaml, so all you need to do is prettify the output

Given a sample.json file of:
```json
{"cat": "meow"}
```
then
```bash
yq -P '.' sample.json
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
yq -P '.' sample.json
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
yq -o=json '.' sample.yml
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
yq -o=json -I=0 '.' sample.yml
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
yq -o=json '.' sample.yml
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
yq -o=json '.' sample.yml
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
yq -o=json -I=0 '.things[]' sample.yml
```
will output
```json
{"stuff":"cool"}
{"whatever":"cat"}
```

## Roundtrip NDJSON
Unfortunately the json encoder strips leading spaces of values.

Given a sample.json file of:
```json
{"this": "is a multidoc json file"}
{"each": ["line is a valid json document"]}
{"a number": 4}

```
then
```bash
yq -p=json -o=json -I=0 sample.json
```
will output
```yaml
{"this":"is a multidoc json file"}
{"each":["line is a valid json document"]}
{"a number":4}
```

## Roundtrip multi-document JSON
The NDJSON parser can also handle multiple multi-line json documents in a single file!

Given a sample.json file of:
```json
{
	"this": "is a multidoc json file"
}
{
	"it": [
		"has",
		"consecutive",
		"json documents"
	]
}
{
	"a number": 4
}

```
then
```bash
yq -p=json -o=json -I=2 sample.json
```
will output
```yaml
{
  "this": "is a multidoc json file"
}
{
  "it": [
    "has",
    "consecutive",
    "json documents"
  ]
}
{
  "a number": 4
}
```

## Update a specific document in a multi-document json
Documents are indexed by the `documentIndex` or `di` operator.

Given a sample.json file of:
```json
{"this": "is a multidoc json file"}
{"each": ["line is a valid json document"]}
{"a number": 4}

```
then
```bash
yq -p=json -o=json -I=0 '(select(di == 1) | .each ) += "cool"' sample.json
```
will output
```yaml
{"this":"is a multidoc json file"}
{"each":["line is a valid json document","cool"]}
{"a number":4}
```

## Find and update a specific document in a multi-document json
Use expressions as you normally would.

Given a sample.json file of:
```json
{"this": "is a multidoc json file"}
{"each": ["line is a valid json document"]}
{"a number": 4}

```
then
```bash
yq -p=json -o=json -I=0 '(select(has("each")) | .each ) += "cool"' sample.json
```
will output
```yaml
{"this":"is a multidoc json file"}
{"each":["line is a valid json document","cool"]}
{"a number":4}
```

## Decode NDJSON
Given a sample.json file of:
```json
{"this": "is a multidoc json file"}
{"each": ["line is a valid json document"]}
{"a number": 4}

```
then
```bash
yq -p=json sample.json
```
will output
```yaml
this: is a multidoc json file
---
each:
    - line is a valid json document
---
a number: 4
```

