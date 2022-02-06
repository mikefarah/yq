# Tag

The tag operator can be used to get or set the tag of nodes (e.g. `!!str`, `!!int`, `!!bool`).

{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Get tag
Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
f: []
```
then
```bash
yq '.. | tag' sample.yml
```
will output
```yaml
!!map
!!str
!!int
!!float
!!bool
!!seq
```

## Set custom tag
Given a sample.yml file of:
```yaml
a: str
```
then
```bash
yq '.a tag = "!!mikefarah"' sample.yml
```
will output
```yaml
a: !!mikefarah str
```

## Find numbers and convert them to strings
Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
```
then
```bash
yq '(.. | select(tag == "!!int")) tag= "!!str"' sample.yml
```
will output
```yaml
a: cat
b: "5"
c: 3.2
e: true
```

