This is a boolean operator and will return `true` when given a `false` value (including null), and `false` otherwise.
## Examples
### Not true is false
Running
```bash
yq eval --null-input 'true | not'
```
will output
```yaml
false
```

### Not false is true
Running
```bash
yq eval --null-input 'false | not'
```
will output
```yaml
true
```

### String values considered to be true
Running
```bash
yq eval --null-input '"cat" | not'
```
will output
```yaml
false
```

### Empty string value considered to be true
Running
```bash
yq eval --null-input '"" | not'
```
will output
```yaml
false
```

### Numbers are considered to be true
Running
```bash
yq eval --null-input '1 | not'
```
will output
```yaml
false
```

### Zero is considered to be true
Running
```bash
yq eval --null-input '0 | not'
```
will output
```yaml
false
```

### Null is considered to be false
Running
```bash
yq eval --null-input '~ | not'
```
will output
```yaml
true
```

