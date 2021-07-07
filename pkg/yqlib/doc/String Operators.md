# String Operators

## Match with names capture groups
Given a sample.yml file of:
```yaml
a cat
```
then
```bash
yq eval 'match("c(?P<cool>.t)")' sample.yml
```
will output
```yaml
string: cat
offset: 2
length: 3
captures:
  - string: at
    offset: 3
    length: 2
    name: cool
```

