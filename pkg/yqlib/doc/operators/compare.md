# Compare Operators

Comparison operators (`>`, `>=`, `<`, `<=`) can be used for comparing scalar values of the same time.

The following types are currently supported:

- numbers
- strings
- datetimes

## Related Operators

- equals / not equals (`==`, `!=`) operators [here])(https://mikefarah.gitbook.io/yq/operators/equals)
- boolean operators (`and`, `or`, `any` etc) [here](https://mikefarah.gitbook.io/yq/operators/boolean-operators)
- select operator [here](https://mikefarah.gitbook.io/yq/operators/select)

{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Compare numbers (>)
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

## Compare strings
Compares strings by their bytecode.

Given a sample.yml file of:
```yaml
a: zoo
b: apple
```
then
```bash
yq '.a > .b' sample.yml
```
will output
```yaml
true
```

## Compare date times
You can compare date times. Assumes RFC3339 date time format, see [date-time operators](https://mikefarah.gitbook.io/yq/operators/date-time-operators) for more information.

Given a sample.yml file of:
```yaml
a: 2021-01-01T03:10:00Z
b: 2020-01-01T03:10:00Z
```
then
```bash
yq '.a > .b' sample.yml
```
will output
```yaml
true
```

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

