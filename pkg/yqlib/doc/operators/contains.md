# Contains

This returns `true` if the context contains the passed in parameter, and false otherwise.

{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
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

