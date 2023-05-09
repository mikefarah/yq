# Subtract

You can use subtract to subtract numbers as well as remove elements from an array.

## Array subtraction
Running
```bash
yq --null-input '[1,2] - [2,3]'
```
will output
```yaml
- 1
```

## Array subtraction with nested array
Running
```bash
yq --null-input '[[1], 1, 2] - [[1], 3]'
```
will output
```yaml
- 1
- 2
```

## Array subtraction with nested object
Note that order of the keys does not matter

Given a sample.yml file of:
```yaml
[{a: b, c: d}, {a: b}]
```
then
```bash
yq '. - [{"c": "d", "a": "b"}]' sample.yml
```
will output
```yaml
[{a: b}]
```

## Number subtraction - float
If the lhs or rhs are floats then the expression will be calculated with floats.

Given a sample.yml file of:
```yaml
{a: 3, b: 4.5}
```
then
```bash
yq '.a = .a - .b' sample.yml
```
will output
```yaml
{a: -1.5, b: 4.5}
```

## Number subtraction - int
If both the lhs and rhs are ints then the expression will be calculated with ints.

Given a sample.yml file of:
```yaml
{a: 3, b: 4}
```
then
```bash
yq '.a = .a - .b' sample.yml
```
will output
```yaml
{a: -1, b: 4}
```

## Decrement numbers
Given a sample.yml file of:
```yaml
{a: 3, b: 5}
```
then
```bash
yq '.[] -= 1' sample.yml
```
will output
```yaml
{a: 2, b: 4}
```

## Date subtraction
You can subtract durations from dates. Assumes RFC3339 date time format, see [date-time operators](https://mikefarah.gitbook.io/yq/operators/date-time-operators) for more information.

Given a sample.yml file of:
```yaml
a: 2021-01-01T03:10:00Z
```
then
```bash
yq '.a -= "3h10m"' sample.yml
```
will output
```yaml
a: 2021-01-01T00:00:00Z
```

## Date subtraction - custom format
Use with_dtf to specify your datetime format. See [date-time operators](https://mikefarah.gitbook.io/yq/operators/date-time-operators) for more information.

Given a sample.yml file of:
```yaml
a: Saturday, 15-Dec-01 at 6:00AM GMT
```
then
```bash
yq 'with_dtf("Monday, 02-Jan-06 at 3:04PM MST", .a -= "3h1m")' sample.yml
```
will output
```yaml
a: Saturday, 15-Dec-01 at 2:59AM GMT
```

## Custom types: that are really numbers
When custom tags are encountered, yq will try to decode the underlying type.

Given a sample.yml file of:
```yaml
a: !horse 2
b: !goat 1
```
then
```bash
yq '.a -= .b' sample.yml
```
will output
```yaml
a: !horse 1
b: !goat 1
```

