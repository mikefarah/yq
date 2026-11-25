# Empty

Returns no results, like jq's `empty`. Combined with the alternative operator (`//`) it is a convenient way of
printing nothing when a value is missing, instead of `null`.

## Produce no results
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq 'empty' sample.yml
```
will output
```yaml
```

## Fall back to no output instead of null
Handy in shell scripts, where a missing value should result in an empty string rather than the text `null`.

Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq '.b // empty' sample.yml
```
will output
```yaml
```

## Fall back to no output, value present
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq '.a // empty' sample.yml
```
will output
```yaml
cat
```

## Empty is dropped from collected results
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq '[.a, empty, .a]' sample.yml
```
will output
```yaml
- cat
- cat
```

