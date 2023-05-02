# Keys

Use the `keys` operator to return map keys or array indices. 

## Check node is a key
Given a sample.yml file of:
```yaml
a: frog
```
then
```bash
yq '[... | { "p": path | join("."), "isKey": is_key, "tag": tag }]' sample.yml
```
will output
```yaml
- p: ""
  isKey: false
  tag: '!!map'
- p: a
  isKey: true
  tag: null
  '!!str': null
- p: a
  isKey: false
  tag: '!!str'
```

