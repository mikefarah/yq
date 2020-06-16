# Working with JSON

## Yaml to Json

To convert output to json, use the `--tojson` \(or `-j`\) flag. This is supported by all commands. You can change the json output format by using the [pretty print](output-format.md#pretty-print) or [indent](output-format.md#indent) flags. _Note that due to the implementation of the JSON marshaller in GO, object keys will be sorted on output \(_[_https://golang.org/pkg/encoding/json/\#Marshal_](https://golang.org/pkg/encoding/json/#Marshal)_\)._

Given a sample.yaml file of:

```yaml
b:
  c: 2
```

then

```bash
yq r -j sample.yaml
```

will output

```javascript
{"b":{"c":2}}
```

To format the json:

```yaml
yq r --prettyPrint -j sample.yaml
```

will yield

```yaml
{
  "b": {
    "c": 2
  }
}
```

### Multiple matches

Each matching yaml node will be converted to json and printed out on a separate line. The [prettyPrint](output-format.md#pretty-print) and [indent](output-format.md#indent) flags will still work too. ****

Given a sample.yaml file of:

```yaml
bob:
  c: 2
bab:
  c: 5
```

then

```bash
yq r -j sample.yaml b*
```

will output

```javascript
{"c":2}
{"c":5}
```

## Json to Yaml

To read in json, just pass in a json file instead of yaml, it will just work - as json is a subset of yaml. However, you will probably want to [pretty print the output](output-format.md#pretty-print) to look more like an idiomatic yaml document.

e.g given a json file

```javascript
{"a":"Easy! as one two three","b":{"c":2,"d":[3,4]}}
```

then

```bash
yq r --prettyPrint sample.json
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

