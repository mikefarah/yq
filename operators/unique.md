# Unique

This is used to filter out duplicated items in an array. Note that the original order of the array is maintained.


{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Unique array of scalars (string/numbers)
Note that unique maintains the original order of the array.

Given a sample.yml file of:
```yaml
- 2
- 1
- 3
- 2
```
then
```bash
yq 'unique' sample.yml
```
will output
```yaml
- 2
- 1
- 3
```

## Unique nulls
Unique works on the node value, so it considers different representations of nulls to be different

Given a sample.yml file of:
```yaml
- ~
- null
- ~
- null
```
then
```bash
yq 'unique' sample.yml
```
will output
```yaml
- ~
- null
```

## Unique all nulls
Run against the node tag to unique all the nulls

Given a sample.yml file of:
```yaml
- ~
- null
- ~
- null
```
then
```bash
yq 'unique_by(tag)' sample.yml
```
will output
```yaml
- ~
```

## Unique array object fields
Given a sample.yml file of:
```yaml
- name: harry
  pet: cat
- name: billy
  pet: dog
- name: harry
  pet: dog
```
then
```bash
yq 'unique_by(.name)' sample.yml
```
will output
```yaml
- name: harry
  pet: cat
- name: billy
  pet: dog
```

