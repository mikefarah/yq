# Filter

Filters an array (or map values) by the expression given. Equivalent to doing `map(select(exp))`.


## Filter array
Given a sample.yml file of:
```yaml
[1, 2, 3]
```
then
```bash
yq 'filter(. < 3)' sample.yml
```
will output
```yaml
[1, 2]
```

## Filter map values
Given a sample.yml file of:
```yaml
{c: {things: cool, frog: yes}, d: {things: hot, frog: false}}
```
then
```bash
yq 'filter(.things == "cool")' sample.yml
```
will output
```yaml
[{things: cool, frog: yes}]
```

