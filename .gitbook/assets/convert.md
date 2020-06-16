## Yaml to Json
To convert output to json, use the --tojson (or -j) flag. This is supported by all commands.

Each matching yaml node will be converted to json and printed out on a separate line.

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
```json
{"b":{"c":2}}
```

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
```json
{"c":2}
{"c":5}
```

## Json to Yaml
To read in json, just pass in a json file instead of yaml, it will just work :)

e.g given a json file

```json
{"a":"Easy! as one two three","b":{"c":2,"d":[3,4]}}
```
then
```bash
yq r sample.json
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

