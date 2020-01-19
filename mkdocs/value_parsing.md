This describes how values are parsed from the CLI to commands that create/update yaml (e.g. new/write).

`yq` attempts to parse values intelligently, e.g. when a number is passed it - it will assume it's a number as opposed to a string. `yq` will not alter the representation of what you give. So if you pass '03.0' in, it will assume it's a number and keep the value formatted as it was passed in, that is '03.0'.

The `--tag` flag can be used to override the tag type to force particular tags.


## Default behaviour

### Integers
*Given*
```bash
yq new key 3
```

results in

```yaml
key: 3
```

*Given a formatted number*

```bash
yq new key 03
```

results in

```yaml
key: 03
```
`yq` keeps the number formatted as it was passed in.

### Float
*Given*
```bash
yq new key "3.1"
```

results in

```yaml
key: 3.1
```
Note that quoting the number does not make a difference.

*Given a formatted decimal number*

```bash
yq new key 03.0
```

results in 

```yaml
key: 03.0
```
`yq` keeps the number formatted as it was passed in

### Booleans
```bash
yq new key true
```

results in

```yaml
key: true
```

### Nulls
```bash
yq new key null
```

results in

```yaml
key: null
```

```bash
yq new key '~'
```

results in

```yaml
key: ~
```

```bash
yq new key ''
```

results in

```yaml
key:
```

### Strings
```bash
yq new key whatever
```

results in

```yaml
key: whatever
```

```bash
yq new key ' whatever '
```

results in

```yaml
key: ' whatever '
```

## Using the tag field to override

Previous versions of yq required double quoting to force values to be strings, this no longer works - instead use the --tag flag.



## Casting booleans
```bash
yq new --tag '!!str' key true
```

results in

```yaml
key: 'true'
```

## Casting nulls
```bash
yq new --tag '!!str' key null
```

results in

```yaml
key: 'null'
```

## Custom types
```bash
yq new --tag '!!farah' key gold
```

results in

```yaml
key: !!farah gold
```