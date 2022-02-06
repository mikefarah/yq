# Collect into Array

This creates an array using the expression between the square brackets.


{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
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

