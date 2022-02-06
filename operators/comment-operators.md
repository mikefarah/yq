# Comment Operators

Use these comment operators to set or retrieve comments.

Like the `=` and `|=` assign operators, the same syntax applies when updating comments:

### plain form: `=`
This will assign the LHS nodes comments to the expression on the RHS. The RHS is run against the matching nodes in the pipeline

### relative form: `|=` 
Similar to the plain form, however the RHS evaluates against each matching LHS node! This is useful if you want to set the comments as a relative expression of the node, for instance its value or path.

{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Set line comment
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq '.a lineComment="single"' sample.yml
```
will output
```yaml
a: cat # single
```

## Use update assign to perform relative updates
Given a sample.yml file of:
```yaml
a: cat
b: dog
```
then
```bash
yq '.. lineComment |= .' sample.yml
```
will output
```yaml
a: cat # cat
b: dog # dog
```

## Set head comment
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq '. headComment="single"' sample.yml
```
will output
```yaml
# single

a: cat
```

## Set foot comment, using an expression
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq '. footComment=.a' sample.yml
```
will output
```yaml
a: cat

# cat
```

## Remove comment
Given a sample.yml file of:
```yaml
a: cat # comment
b: dog # leave this
```
then
```bash
yq '.a lineComment=""' sample.yml
```
will output
```yaml
a: cat
b: dog # leave this
```

## Remove (strip) all comments
Note the use of `...` to ensure key nodes are included.

Given a sample.yml file of:
```yaml
a: cat # comment
# great
b: # key comment
```
then
```bash
yq '... comments=""' sample.yml
```
will output
```yaml
a: cat
b:
```

## Get line comment
Given a sample.yml file of:
```yaml
a: cat # meow
```
then
```bash
yq '.a | lineComment' sample.yml
```
will output
```yaml
meow
```

## Get head comment
Given a sample.yml file of:
```yaml
# welcome!

a: cat # meow

# have a great day
```
then
```bash
yq '. | headComment' sample.yml
```
will output
```yaml
welcome!
```

## Head comment with document split
Given a sample.yml file of:
```yaml
# welcome!
---
# bob
a: cat # meow

# have a great day
```
then
```bash
yq 'headComment' sample.yml
```
will output
```yaml
welcome!
bob
```

## Get foot comment
Given a sample.yml file of:
```yaml
# welcome!

a: cat # meow

# have a great day
# no really
```
then
```bash
yq '. | footComment' sample.yml
```
will output
```yaml
have a great day
no really
```

