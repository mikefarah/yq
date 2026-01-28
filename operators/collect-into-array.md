# Collect into Array

This creates an array using the expression between the square brackets.

{% hint style="warning" %}

_Note_ the placement of `|` when collecting. These two forms behave differently:

```bash
# Pipe then splat - creates separate context
[.items | .[] | has("id")]

# Splat directly on path - more common pattern
[.items[] | has("id")]
```

{% endhint %}

## Collect empty
Running
```bash
yq --null-input '[]'
```
will output
```yaml
[]
```

## Collect single
Running
```bash
yq --null-input '["cat"]'
```
will output
```yaml
- cat
```

## Collect many
Given a sample.yml file of:
```yaml
a: cat
b: dog
```
then
```bash
yq '[.a, .b]' sample.yml
```
will output
```yaml
- cat
- dog
```

