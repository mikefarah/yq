### Yaml to Json
To convert output to json, use the --tojson (or -j) flag. This can only be used with the read command.

Given a sample.yaml file of:
```yaml
b:
  c: 2
```
then
```bash
yq r -j sample.yaml b.c
```

will output
```json
{"b":{"c":2}}
```

### Json to Yaml
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

