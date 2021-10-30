# Boolean Operators

The `or` and `and` operators take two parameters and return a boolean result.

`not` flips a boolean from true to false, or vice versa.

`any` will return `true` if there are any `true` values in a array sequence, and `all` will return true if _all_ elements in an array are true.

`any_c(condition)` and `all_c(condition)` are like `any` and `all` but they take a condition expression that is used against each element to determine if it's `true`. Note: in `jq` you can simply pass a condition to `any` or `all` and it simply works - `yq` isn't that clever..yet

These are most commonly used with the `select` operator to filter particular nodes.

## `or` example

Running

```bash
yq eval --null-input 'true or false'
```

will output

```yaml
true
```

## `and` example

Running

```bash
yq eval --null-input 'true and false'
```

will output

```yaml
false
```

## Matching nodes with select, equals and or

Given a sample.yml file of:

```yaml
- a: bird
  b: dog
- a: frog
  b: bird
- a: cat
  b: fly
```

then

```bash
yq eval '[.[] | select(.a == "cat" or .b == "dog")]' sample.yml
```

will output

```yaml
- a: bird
  b: dog
- a: cat
  b: fly
```

## `any` returns true if any boolean in a given array is true

Given a sample.yml file of:

```yaml
- false
- true
```

then

```bash
yq eval 'any' sample.yml
```

will output

```yaml
true
```

## `any` returns false for an empty array

Given a sample.yml file of:

```yaml
[]
```

then

```bash
yq eval 'any' sample.yml
```

will output

```yaml
false
```

## `any_c` returns true if any element in the array is true for the given condition.

Given a sample.yml file of:

```yaml
a:
  - rad
  - awesome
b:
  - meh
  - whatever
```

then

```bash
yq eval '.[] |= any_c(. == "awesome")' sample.yml
```

will output

```yaml
a: true
b: false
```

## `all` returns true if all booleans in a given array are true

Given a sample.yml file of:

```yaml
- true
- true
```

then

```bash
yq eval 'all' sample.yml
```

will output

```yaml
true
```

## `all` returns true for an empty array

Given a sample.yml file of:

```yaml
[]
```

then

```bash
yq eval 'all' sample.yml
```

will output

```yaml
true
```

## `all_c` returns true if all elements in the array are true for the given condition.

Given a sample.yml file of:

```yaml
a:
  - rad
  - awesome
b:
  - meh
  - 12
```

then

```bash
yq eval '.[] |= all_c(tag == "!!str")' sample.yml
```

will output

```yaml
a: true
b: false
```

## Not true is false

Running

```bash
yq eval --null-input 'true | not'
```

will output

```yaml
false
```

## Not false is true

Running

```bash
yq eval --null-input 'false | not'
```

will output

```yaml
true
```

## String values considered to be true

Running

```bash
yq eval --null-input '"cat" | not'
```

will output

```yaml
false
```

## Empty string value considered to be true

Running

```bash
yq eval --null-input '"" | not'
```

will output

```yaml
false
```

## Numbers are considered to be true

Running

```bash
yq eval --null-input '1 | not'
```

will output

```yaml
false
```

## Zero is considered to be true

Running

```bash
yq eval --null-input '0 | not'
```

will output

```yaml
false
```

## Null is considered to be false

Running

```bash
yq eval --null-input '~ | not'
```

will output

```yaml
true
```
