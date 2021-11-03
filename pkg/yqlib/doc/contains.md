# Contains

This returns `true` if the context contains the passed in parameter, and false otherwise.

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
yq eval 'contains(["baz", "bar"])' sample.yml
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
yq eval 'contains({"bar": [{"barp": 12}]})' sample.yml
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
yq eval 'contains({"foo": 12, "bar": [{"barp": 15}]})' sample.yml
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
yq eval 'contains("bar")' sample.yml
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
yq eval 'contains("meow")' sample.yml
```
will output
```yaml
true
```

