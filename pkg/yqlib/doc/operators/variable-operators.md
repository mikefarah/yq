# Variable Operators

Like the `jq` equivalents, variables are sometimes required for the more complex expressions (or swapping values between fields).

Note that there is also an additional `ref` operator that holds a reference (instead of a copy) of the path, allowing you to make multiple changes to the same path.

## Single value variable
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq eval '.a as $foo | $foo' sample.yml
```
will output
```yaml
cat
```

## Multi value variable
Given a sample.yml file of:
```yaml
- cat
- dog
```
then
```bash
yq eval '.[] as $foo | $foo' sample.yml
```
will output
```yaml
cat
dog
```

## Using variables as a lookup
Example taken from [jq](https://stedolan.github.io/jq/manual/#Variable/SymbolicBindingOperator:...as$identifier|...)

Given a sample.yml file of:
```yaml
"posts":
  - "title": Frist psot
    "author": anon
  - "title": A well-written article
    "author": person1
"realnames":
  "anon": Anonymous Coward
  "person1": Person McPherson
```
then
```bash
yq eval '.realnames as $names | .posts[] | {"title":.title, "author": $names[.author]}' sample.yml
```
will output
```yaml
title: Frist psot
author: Anonymous Coward
title: A well-written article
author: Person McPherson
```

## Using variables to swap values
Given a sample.yml file of:
```yaml
a: a_value
b: b_value
```
then
```bash
yq eval '.a as $x  | .b as $y | .b = $x | .a = $y' sample.yml
```
will output
```yaml
a: b_value
b: a_value
```

## Use ref to reference a path repeatedly
Note: You may find the `with` operator more useful.

Given a sample.yml file of:
```yaml
a:
  b: thing
  c: something
```
then
```bash
yq eval '.a.b ref $x | $x = "new" | $x style="double"' sample.yml
```
will output
```yaml
a:
  b: "new"
  c: something
```

