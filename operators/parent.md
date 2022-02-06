# Parent

Parent simply returns the parent nodes of the matching nodes.

{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Simple example
Given a sample.yml file of:
```yaml
a:
  nested: cat
```
then
```bash
yq '.a.nested | parent' sample.yml
```
will output
```yaml
nested: cat
```

## Parent of nested matches
Given a sample.yml file of:
```yaml
a:
  fruit: apple
  name: bob
b:
  fruit: banana
  name: sam
```
then
```bash
yq '.. | select(. == "banana") | parent' sample.yml
```
will output
```yaml
fruit: banana
name: sam
```

## No parent
Given a sample.yml file of:
```yaml
{}
```
then
```bash
yq 'parent' sample.yml
```
will output
```yaml
```

