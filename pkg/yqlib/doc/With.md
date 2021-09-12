Use the `with` operator to conveniently make multiple updates to a deeply nested path.

## Update and style
Given a sample.yml file of:
```yaml
a:
  deeply:
    nested: value
```
then
```bash
yq eval 'with(.a.deeply.nested ; . = "newValue" | . style="single")' sample.yml
```
will output
```yaml
a:
  deeply:
    nested: 'newValue'
```

