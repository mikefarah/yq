# String Operators

## Match string
Given a sample.yml file of:
```yaml
cat
```
then
```bash
yq eval 'match("at")' sample.yml
```
will output
```yaml
string: at
offset: 1
length: 2
captures: []
```

## Match string, case insensitive
Given a sample.yml file of:
```yaml
cAt
```
then
```bash
yq eval 'match("(?i)at")' sample.yml
```
will output
```yaml
string: At
offset: 1
length: 2
captures: []
```

## Match with capture groups
Given a sample.yml file of:
```yaml
a cat
```
then
```bash
yq eval 'match("c(.t)")' sample.yml
```
will output
```yaml
string: cat
offset: 2
length: 3
captures:
  - string: at
    offset: 3
    length: 2
```

