Encode operators will take the piped in object structure and encode it as a string in the desired format.
## Encode value as yaml string
Given a sample.yml file of:
```yaml
a:
  cool: thing
```
then
```bash
yq eval '.b = (.a | to_yaml)' sample.yml
```
will output
```yaml
a:
  cool: thing
b: |
  cool: thing
```

## Encode value as yaml string, using toyaml
Does the same thing as to_yaml, matching jq naming convention.

Given a sample.yml file of:
```yaml
a:
  cool: thing
```
then
```bash
yq eval '.b = (.a | to_yaml)' sample.yml
```
will output
```yaml
a:
  cool: thing
b: |
  cool: thing
```

## Encode value as json string
Given a sample.yml file of:
```yaml
a:
  cool: thing
```
then
```bash
yq eval '.b = (.a | to_json)' sample.yml
```
will output
```yaml
a:
  cool: thing
b: |
  {
    "cool": "thing"
  }
```

## Encode value as props string
Given a sample.yml file of:
```yaml
a:
  cool: thing
```
then
```bash
yq eval '.b = (.a | to_props)' sample.yml
```
will output
```yaml
a:
  cool: thing
b: |
  cool = thing
```

