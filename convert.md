# Convert

### Yaml to Json[¶](convert.md#yaml-to-json) <a id="yaml-to-json"></a>

To convert output to json, use the --tojson \(or -j\) flag. This can only be used with the read command.

Given a sample.yaml file of:

```text
b:
  c: 2
```

then

```text
yq r -j sample.yaml b.c
```

will output

```text
{"b":{"c":2}}
```

### Json to Yaml[¶](convert.md#json-to-yaml) <a id="json-to-yaml"></a>

To read in json, just pass in a json file instead of yaml, it will just work :\)

e.g given a json file

```text
{"a":"Easy! as one two three","b":{"c":2,"d":[3,4]}}
```

then

```text
yq r sample.json
```

will output

```text
a: Easy! as one two three
b:
  c: 2
  d:
  - 3
  - 4
```

