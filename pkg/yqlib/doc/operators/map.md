# Map

Maps values of an array. Use `map_values` to map values of an object.

## Map array
Given a sample.yml file of:
```yaml
- 1
- 2
- 3
```
then
```bash
yq eval 'map(. + 1)' sample.yml
```
will output
```yaml
- 2
- 3
- 4
```

## Map object values
Given a sample.yml file of:
```yaml
a: 1
b: 2
c: 3
```
then
```bash
yq eval 'map_values(. + 1)' sample.yml
```
will output
```yaml
a: 2
b: 3
c: 4
```

