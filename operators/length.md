# Length

Returns the lengths of the nodes. Length is defined according to the type of the node.

{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## String length
returns length of string

Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq '.a | length' sample.yml
```
will output
```yaml
3
```

## null length
Given a sample.yml file of:
```yaml
a: null
```
then
```bash
yq '.a | length' sample.yml
```
will output
```yaml
0
```

## Map length
returns number of entries

Given a sample.yml file of:
```yaml
a: cat
c: dog
```
then
```bash
yq 'length' sample.yml
```
will output
```yaml
2
```

## Array length
returns number of elements

Given a sample.yml file of:
```yaml
- 2
- 4
- 6
- 8
```
then
```bash
yq 'length' sample.yml
```
will output
```yaml
4
```

