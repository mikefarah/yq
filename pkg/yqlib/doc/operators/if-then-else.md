# If Then Else

`if A then B else C end` evaluates `A` against each input, and returns `B` if the result is truthy, otherwise `C`. Like `select`, `null` and `false` are falsy, and everything else is truthy.

`elif` can be used to chain conditions, and `else` is optional - if it is left out and the condition is falsy the input is returned unchanged (like jq 1.7).

## Related Operators

- select operator [here](https://mikefarah.gitbook.io/yq/operators/select)
- alternative (`//`) operator [here](https://mikefarah.gitbook.io/yq/operators/alternative-default-value)
- boolean operators (`and`, `or`, `any` etc) [here](https://mikefarah.gitbook.io/yq/operators/boolean-operators)

## Basic if-then-else
The condition is evaluated against each input; `null` and `false` are falsy, everything else is truthy.

Given a sample.yml file of:
```yaml
- 1
- 2
- 3
```
then
```bash
yq '.[] |= if . > 1 then "big" else "small" end' sample.yml
```
will output
```yaml
- small
- big
- big
```

## elif
Given a sample.yml file of:
```yaml
- 1
- 2
- 3
```
then
```bash
yq '.[] | if . == 1 then "one" elif . == 2 then "two" else "many" end' sample.yml
```
will output
```yaml
one
two
many
```

## else is optional
Like jq, if there is no `else` and the condition is falsy, the input is returned unchanged.

Given a sample.yml file of:
```yaml
- 1
- 2
- 3
```
then
```bash
yq '.[] |= if . == 2 then "two" end' sample.yml
```
will output
```yaml
- 1
- two
- 3
```

## Update matching branch
The branches can be used to choose which path to update.

Given a sample.yml file of:
```yaml
enabled: true
a: 1
b: 2
```
then
```bash
yq '(if .enabled then .a else .b end) = 10' sample.yml
```
will output
```yaml
enabled: true
a: 10
b: 2
```

## Condition with multiple results
Like jq, each result of the condition produces a result.

Running
```bash
yq --null-input 'if true, false then "yes" else "no" end'
```
will output
```yaml
yes
no
```

## Keywords can still be used as keys
Given a sample.yml file of:
```yaml
if:
  then:
    else:
      end: 1
```
then
```bash
yq 'if .if.then.else.end == 1 then "OK" else "OOPS" end' sample.yml
```
will output
```yaml
OK
```

## Missing keys are falsy
A condition that returns no results, such as a missing key, is treated as `null` - consistent with `and`, `or` and `//`.

Given a sample.yml file of:
```yaml
a: 1
```
then
```bash
yq 'if .b then "has b" else "no b" end' sample.yml
```
will output
```yaml
no b
```

