# Equals

This is a boolean operator that will return `true` if the LHS is equal to the RHS and `false` otherwise.

```
.a == .b
```

It is most often used with the select operator to find particular nodes:

```
select(.a == .b)
```

## Match string
Given a sample.yml file of:
```yaml
- cat
- goat
- dog
```
then
```bash
yq eval '.[] | (. == "*at")' sample.yml
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
yq eval '.[] | (. != "*at")' sample.yml
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
yq eval '.[] | (. == 4)' sample.yml
```
will output
```yaml
false
true
false
```

## Dont match number
Given a sample.yml file of:
```yaml
- 3
- 4
- 5
```
then
```bash
yq eval '.[] | (. != 4)' sample.yml
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
yq eval --null-input 'null == ~'
```
will output
```yaml
true
```

## Non exisitant key doesn't equal a value
Given a sample.yml file of:
```yaml
a: frog
```
then
```bash
yq eval 'select(.b != "thing")' sample.yml
```
will output
```yaml
a: frog
```

## Two non existant keys are equal
Given a sample.yml file of:
```yaml
a: frog
```
then
```bash
yq eval 'select(.b == .c)' sample.yml
```
will output
```yaml
a: frog
```

