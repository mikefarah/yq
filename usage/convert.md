# Working with JSON

## Yaml to Json

To convert output to json, use the `--output-format=json` (or `-o=j`) flag. You can change the json output format by using the [indent](output-format.md#indent) flag.&#x20;

Given a sample.yaml file of:

```yaml
b:
  c: 2
```

then

```bash
yq eval -o=j sample.yaml
```

will output

```javascript
{
  "b": {
    "c": 2
  }
}
```

To format the json:

```yaml
yq eval -o=j -I=0 sample.yaml
```

will yield

```yaml
{"b":{"c":2}}
```

### Multiple matches

Each matching yaml node will be converted to json and printed out as a separate json doc. You may want to set the [indent](output-format.md#indent) flags to 0 if you want a json doc per line.

Given a sample.yaml file of:

```yaml
bob:
  c: 2
bab:
  c: 5
```

then

```bash
yq eval -o=j '.b*' sample.yaml
```

will output

```javascript
{
  "c": 2
}
{
  "c": 5
}
```

## Json to Yaml

To read in json, just pass in a json file instead of yaml, it will just work - as json is a subset of yaml. However, you will probably want to use the [Style Operator](broken-reference) or `--prettyPrint/-P` flag to make look more like an idiomatic yaml document. This can be done by resetting the style of all elements.

e.g given a json file

```javascript
{"a":"Easy! as one two three","b":{"c":2,"d":[3,4]}}
```

then

```bash
yq eval -P sample.json
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
