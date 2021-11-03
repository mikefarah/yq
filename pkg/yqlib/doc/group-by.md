# Group By

This is used to group items in an array by an expression.

## Group by field
Given a sample.yml file of:
```yaml
- foo: 1
  bar: 10
- foo: 3
  bar: 100
- foo: 1
  bar: 1
```
then
```bash
yq eval 'group_by(.foo)' sample.yml
```
will output
```yaml
- - foo: 1
    bar: 10
  - foo: 1
    bar: 1
- - foo: 3
    bar: 100
```

## Group by field, with nuls
Given a sample.yml file of:
```yaml
- cat: dog
- foo: 1
  bar: 10
- foo: 3
  bar: 100
- no: foo for you
- foo: 1
  bar: 1
```
then
```bash
yq eval 'group_by(.foo)' sample.yml
```
will output
```yaml
- - cat: dog
  - no: foo for you
- - foo: 1
    bar: 10
  - foo: 1
    bar: 1
- - foo: 3
    bar: 100
```

