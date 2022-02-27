# Reverse

Reverses the order of the items in an array 

{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Reverse
Given a sample.yml file of:
```yaml
- 1
- 2
- 3
```
then
```bash
yq 'reverse' sample.yml
```
will output
```yaml
- 3
- 2
- 1
```

## Sort descending by string field
Use sort with reverse to sort in descending order.

Given a sample.yml file of:
```yaml
- a: banana
- a: cat
- a: apple
```
then
```bash
yq 'sort_by(.a) | reverse' sample.yml
```
will output
```yaml
- a: cat
- a: banana
- a: apple
```

