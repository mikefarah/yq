
## Get anchor
Given a sample.yml file of:
```yaml
a: &billyBob cat
```
then
```bash
yq eval '.a | anchor' sample.yml
```
will output
```yaml
billyBob
```

## Set anchor name
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq eval '.a anchor = "foobar"' sample.yml
```
will output
```yaml
a: &foobar cat
```

