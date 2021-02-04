For more complex scenarios, variables can be used to hold values of expression to be used in other expressions.

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

