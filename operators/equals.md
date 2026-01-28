# Equals / Not Equals

This is a boolean operator that will return `true` if the LHS is equal to the RHS and `false` otherwise.

```
.a == .b
```

It is most often used with the select operator to find particular nodes:

```
select(.a == .b)
```

The not equals `!=` operator returns `false` if the LHS is equal to the RHS.

{% hint style="info" %}

_Note_ that the `!=` operator contains `!` which can trigger bash history expansion in interactive shells:

```bash
# May fail with: bash: !": event not found
yq '.status != "healthy"' file.yaml
```

Workarounds:
```bash
# Use $'...' quoting
yq $'.status != "healthy"' file.yaml

# Or use | not instead of !=
yq '.status == "healthy" | not' file.yaml
```

{% endhint %}

## Related Operators

- comparison (`>=`, `<` etc) operators [here](https://mikefarah.gitbook.io/yq/operators/compare)
- boolean operators (`and`, `or`, `any` etc) [here](https://mikefarah.gitbook.io/yq/operators/boolean-operators)
- select operator [here](https://mikefarah.gitbook.io/yq/operators/select)


## Match string
Given a sample.yml file of:
```yaml
- cat
- goat
- dog
```
then
```bash
yq '.[] | (. == "*at")' sample.yml
```
will output
```yaml
true
true
false
```

## Don't match string
Given a sample.yml file of:
```yaml
- cat
- goat
- dog
```
then
```bash
yq '.[] | (. != "*at")' sample.yml
```
will output
```yaml
false
false
true
```

## Match number
Given a sample.yml file of:
```yaml
- 3
- 4
- 5
```
then
```bash
yq '.[] | (. == 4)' sample.yml
```
will output
```yaml
false
true
false
```

## Don't match number
Given a sample.yml file of:
```yaml
- 3
- 4
- 5
```
then
```bash
yq '.[] | (. != 4)' sample.yml
```
will output
```yaml
true
false
true
```

## Match nulls
Running
```bash
yq --null-input 'null == ~'
```
will output
```yaml
true
```

## Non existent key doesn't equal a value
Given a sample.yml file of:
```yaml
a: frog
```
then
```bash
yq 'select(.b != "thing")' sample.yml
```
will output
```yaml
a: frog
```

## Two non existent keys are equal
Given a sample.yml file of:
```yaml
a: frog
```
then
```bash
yq 'select(.b == .c)' sample.yml
```
will output
```yaml
a: frog
```

