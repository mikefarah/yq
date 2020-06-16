---
description: Returns matching nodes/values of a path expression for a given yaml document
---

# Read

```text
yq r <yaml_file|json_file> <path_expression>
```

Returns the matching nodes of the path expression for the given yaml file \(or STDIN\).

See docs for [path expression](../usage/path-expressions.md) for more details.

## Basic

Given a sample.yaml file of:

```yaml
b:
  c: 2
```

then

```bash
yq r sample.yaml b.c
```

will output the value of '2'.

## From Stdin

Given a sample.yaml file of:

```bash
cat sample.yaml | yq r - b.c
```

will output the value of '2'.

## Default values

Using the `--defaultValue/-D` flag you can specify a default value to be printed when no matching nodes are found for an expression

```text
yq r sample.yaml --defaultValue frog path.not.there
```

will yield \(assuming `path.not.there` does not match any nodes\):

```text
frog
```

## Printing matching paths

By default, yq will only print the value of the path expression for the yaml document. By specifying `--printMode` or `-p` you can print the matching paths.

```yaml
a:
  thing_a: 
    animal: cat
  other: 
    animal: frog
  thing_b: 
    vehicle: car
```

### Path Only

```bash
yq r --printMode p "a.thing*.*"
```

will print

```text
a.thing_a.animal
a.thing_b.vehicle
```

### Path and Value

```bash
yq r --printMode pv "a.thing*.*"
```

will print

```text
a.thing_a.animal: cat
a.thing_b.vehicle: car
```

## Collect results into an array

By default, results are printed out line by line as independent matches. This is handy for both readability as well as piping into tools like `xargs`. However, if you would like to collect the matching results into an array then use the `--collect/-C` flag. This is particularly useful with the `length` flag described below.

Given:

```yaml
a:
  thing_a: 
    animal: cat
  other: 
    animal: frog
  thing_b: 
    vehicle: car
```

```text
yq r sample.yaml --collect a.*.animal
```

will print

```text
- cat
- frog
```

## Printing length of the results

Use the `--length/-L` flag to print the length of results. For arrays this will be the number of items, objects the number of entries and scalars the length of the value.

Given

```text
animals:
  - cats
  - dog
  - cheetah
```

```text
yq r sample.yml --length animals
```

will print 

```text
3
```

### Length of filtered results

By default, filtered results are printed _independently_ so you will get the length of each result printed on a separate line:

```text
yq r sample.yaml --length --printMode pv 'animals.(.==c*)'
```

```text
animals.[0]: 4
animals.[2]: 7
```

However, you'll often want to know the count of the filtered results - use the `--collect` flag to collect the results in the array. The length will then return the size of the array. 

```text
yq r sample.yaml --length --collect 'animals.(.==c*)' 
```

```text
2
```

## Anchors and Aliases

The read command will print out the anchors of a document and can also traverse them.

Lets take a look at a simple example file:

```yaml
foo: &foo
  a: 1

foobar: *foo
```

### Printing anchors

```bash
yq r sample.yml foo
```

will print out

```yaml
&foo
a: 1
```

Similarly,

```text
yq r sample.yml foobar
```

prints out

```yaml
*foo
```

### Traversing anchors

To traverse an anchor, we need to either explicitly reference merged in values:

```text
yq r sample.yml foobar.a
```

to get

```text
1
```

or we can use deep splat to get all the values

```bash
yq r sample.yml -p pv "foobar.**"
```

prints

```yaml
foobar.a: 1
```

The same methods work for the `<<: [*blah, *thing]`anchors.

### Exploding Anchors

By default anchors are not exploded \(or expanded/de-referenced\) for viewing, and the yaml is shown as-is. Use the `--explodeAnchors/-X` flag to show the anchor values.

Given sample.yml:

```yaml
foo: &foo
  a: original
  thing: coolasdf
  thirsty: yep

bar: &bar
  b: 2
  thing: coconut
  c: oldbar

foobarList:
  <<: [*foo,*bar]
  c: newbar
```

Then

```text
yq r -X sample.yml foobarList
```

yields

```text
c: newbar
b: 2
thing: coconut
a: original
thirsty: yep
```

Note that yq processes the merge anchor list in reverse order, to ensure that the last items in the list override the preceding.

## Multiple Documents

### Reading from a single document

Given a sample.yaml file of:

```yaml
something: else
---
b:
  c: 2
```

then

```bash
yq r -d1 sample.yaml b.c
```

will output the value of '2'.

### Read from all documents

Reading all documents will return the result as an array. This can be converted to json using the '-j' flag if desired.

Given a sample.yaml file of:

```yaml
name: Fred
age: 22
---
name: Stella
age: 23
---
name: Android
age: 232
```

then

```bash
yq r -d'*' sample.yaml name
```

will output:

```text
Fred
Stella
Android
```

