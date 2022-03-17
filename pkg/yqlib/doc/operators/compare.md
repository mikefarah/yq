
{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Both sides are null: > is false
Running
```bash
yq --null-input '.a > .b'
```
will output
```yaml
false
```

## Both sides are null: >= is true
Running
```bash
yq --null-input '.a >= .b'
```
will output
```yaml
true
```

## One side is null: > is false
Given a sample.yml file of:
```yaml
a: 5
```
then
```bash
yq '.a > .b' sample.yml
```
will output
```yaml
false
```

## One side is null: >= is false
Given a sample.yml file of:
```yaml
a: 5
```
then
```bash
yq '.a >= .b' sample.yml
```
will output
```yaml
false
```

## Compare integers (>)
Given a sample.yml file of:
```yaml
a: 5
b: 4
```
then
```bash
yq '.a > .b' sample.yml
```
will output
```yaml
true
```

## Compare integers (>=)
Given a sample.yml file of:
```yaml
a: 5
b: 4
```
then
```bash
yq '.a >= .b' sample.yml
```
will output
```yaml
true
```

## Compare equal numbers
Given a sample.yml file of:
```yaml
a: 5
b: 5
```
then
```bash
yq '.a > .b' sample.yml
```
will output
```yaml
false
```

## Compare equal numbers (>=)
Given a sample.yml file of:
```yaml
a: 5
b: 5
```
then
```bash
yq '.a >= .b' sample.yml
```
will output
```yaml
true
```

