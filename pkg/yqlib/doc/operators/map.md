# Map

Maps values of an array. Use `map_values` to map values of an object.

{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Map array
Given a sample.yml file of:
```yaml
- 1
- 2
- 3
```
then
```bash
yq 'map(. + 1)' sample.yml
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
yq 'map_values(. + 1)' sample.yml
```
will output
```yaml
a: 2
b: 3
c: 4
```

