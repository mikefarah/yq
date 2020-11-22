Select is used to filter arrays and maps by a boolean expression.
## Select elements from array
Given a sample.yml file of:
```yaml
- cat
- goat
- dog
```
then
```bash
yq eval '.[] | select(. == "*at")' sample.yml
```
will output
```yaml
cat
goat
```

## Select and update matching values in map
Given a sample.yml file of:
```yaml
a:
  things: cat
  bob: goat
  horse: dog
```
then
```bash
yq eval '(.a[] | select(. == "*at")) |= "rabbit"' sample.yml
```
will output
```yaml
a:
  things: rabbit
  bob: rabbit
  horse: dog
```

