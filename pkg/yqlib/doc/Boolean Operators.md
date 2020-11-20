The `or` and `and` operators take two parameters and return a boolean result. These are most commonly used with the `select` operator to filter particular nodes.
## Examples
### OR example
Running
```bash
yq eval --null-input 'true or false'
```
will output
```yaml
true
```

### AND example
Running
```bash
yq eval --null-input 'true and false'
```
will output
```yaml
false
```

### Matching nodes with select, equals and or
Given a sample.yml file of:
```yaml
- a: bird
  b: dog
- a: frog
  b: bird
- a: cat
  b: fly
```
then
```bash
yq eval '[.[] | select(.a == "cat" or .b == "dog")]' sample.yml
```
will output
```yaml
- a: bird
  b: dog
- a: cat
  b: fly
```

