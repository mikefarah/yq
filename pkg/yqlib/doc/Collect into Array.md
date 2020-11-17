# Collect into Array

This creates an array using the expression between the square brackets.


## Examples
### Collect empty
Running
```bash
yq eval --null-input '[]'
```
will output
```yaml
```

### Collect single
Running
```bash
yq eval --null-input '["cat"]'
```
will output
```yaml
- cat
```

### Collect many
Given a sample.yml file of:
```yaml
a: cat
b: dog
```
then
```bash
yq eval '[.a, .b]' sample.yml
```
will output
```yaml
- cat
- dog
```

