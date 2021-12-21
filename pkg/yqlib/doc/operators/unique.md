# Unique

This is used to filter out duplicated items in an array.

## Unique array of scalars (string/numbers)
Given a sample.yml file of:
```yaml
- 1
- 2
- 3
- 2
```
then
```bash
yq eval 'unique' sample.yml
```
will output
```yaml
- 1
- 2
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
yq eval 'unique' sample.yml
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
yq eval 'unique_by(tag)' sample.yml
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
yq eval 'unique_by(.name)' sample.yml
```
will output
```yaml
- name: harry
  pet: cat
- name: billy
  pet: dog
```

