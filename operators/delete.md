# Delete

Deletes matching entries in maps or arrays.

{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Delete entry in map
Given a sample.yml file of:
```yaml
a: cat
b: dog
```
then
```bash
yq 'del(.b)' sample.yml
```
will output
```yaml
a: cat
```

## Delete nested entry in map
Given a sample.yml file of:
```yaml
a:
  a1: fred
  a2: frood
```
then
```bash
yq 'del(.a.a1)' sample.yml
```
will output
```yaml
a:
  a2: frood
```

## Delete entry in array
Given a sample.yml file of:
```yaml
- 1
- 2
- 3
```
then
```bash
yq 'del(.[1])' sample.yml
```
will output
```yaml
- 1
- 3
```

## Delete nested entry in array
Given a sample.yml file of:
```yaml
- a: cat
  b: dog
```
then
```bash
yq 'del(.[0].a)' sample.yml
```
will output
```yaml
- b: dog
```

## Delete no matches
Given a sample.yml file of:
```yaml
a: cat
b: dog
```
then
```bash
yq 'del(.c)' sample.yml
```
will output
```yaml
a: cat
b: dog
```

## Delete matching entries
Given a sample.yml file of:
```yaml
a: cat
b: dog
c: bat
```
then
```bash
yq 'del( .[] | select(. == "*at") )' sample.yml
```
will output
```yaml
b: dog
```

## Recursively delete matching keys
Given a sample.yml file of:
```yaml
a:
  name: frog
  b:
    name: blog
    age: 12
```
then
```bash
yq 'del(.. | select(has("name")).name)' sample.yml
```
will output
```yaml
a:
  b:
    age: 12
```

