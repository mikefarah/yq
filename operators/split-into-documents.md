# Split into Documents

This operator splits all matches into separate documents

{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Split empty
Running
```bash
yq --null-input 'split_doc'
```
will output
```yaml

```

## Split array
Given a sample.yml file of:
```yaml
- a: cat
- b: dog
```
then
```bash
yq '.[] | split_doc' sample.yml
```
will output
```yaml
a: cat
---
b: dog
```

