# Contains

This returns `true` if the context contains the passed in parameter, and false otherwise. For arrays, this will return true if the passed in array is contained within the array. For strings, it will return true if the string is a substring.

{% hint style="warning" %}

_Note_ that, just like jq, when checking if an array of strings `contains` another, this will use `contains` and _not_ equals to check each string. This means an array like `contains(["cat"])` will return true for an array `["cats"]`.

See the "Array has a subset array" example below on how to check for a subset.

{% endhint %}

## Array contains array
Array is equal or subset of

Given a sample.yml file of:
```yaml
- foobar
- foobaz
- blarp
```
then
```bash
yq 'contains(["baz", "bar"])' sample.yml
```
will output
```yaml
true
```

## Array has a subset array
Subtract the superset array from the subset, if there's anything left, it's not a subset

Given a sample.yml file of:
```yaml
- foobar
- foobaz
- blarp
```
then
```bash
yq '["baz", "bar"] - . | length == 0' sample.yml
```
will output
```yaml
false
```

## Object included in array
Given a sample.yml file of:
```yaml
"foo": 12
"bar":
  - 1
  - 2
  - "barp": 12
    "blip": 13
```
then
```bash
yq 'contains({"bar": [{"barp": 12}]})' sample.yml
```
will output
```yaml
true
```

## Object not included in array
Given a sample.yml file of:
```yaml
"foo": 12
"bar":
  - 1
  - 2
  - "barp": 12
    "blip": 13
```
then
```bash
yq 'contains({"foo": 12, "bar": [{"barp": 15}]})' sample.yml
```
will output
```yaml
false
```

## String contains substring
Given a sample.yml file of:
```yaml
foobar
```
then
```bash
yq 'contains("bar")' sample.yml
```
will output
```yaml
true
```

## String equals string
Given a sample.yml file of:
```yaml
meow
```
then
```bash
yq 'contains("meow")' sample.yml
```
will output
```yaml
true
```

