# Add

Add behaves differently according to the type of the LHS:
* arrays: concatenate
* number scalars: arithmetic addition
* string scalars: concatenate
* maps: shallow merge (use the multiply operator (`*`) to deeply merge)

Use `+=` as a relative append assign for things like increment. Note that `.a += .x` is equivalent to running `.a = .a + .x`.


## Concatenate arrays
Given a sample.yml file of:
```yaml
a:
  - 1
  - 2
b:
  - 3
  - 4
```
then
```bash
yq '.a + .b' sample.yml
```
will output
```yaml
- 1
- 2
- 3
- 4
```

## Concatenate to existing array
Note that the styling of `a` is kept.

Given a sample.yml file of:
```yaml
a: [1,2]
b:
  - 3
  - 4
```
then
```bash
yq '.a += .b' sample.yml
```
will output
```yaml
a: [1, 2, 3, 4]
b:
  - 3
  - 4
```

## Concatenate null to array
Given a sample.yml file of:
```yaml
a:
  - 1
  - 2
```
then
```bash
yq '.a + null' sample.yml
```
will output
```yaml
- 1
- 2
```

## Append to existing array
Note that the styling is copied from existing array elements

Given a sample.yml file of:
```yaml
a: ['dog']
```
then
```bash
yq '.a += "cat"' sample.yml
```
will output
```yaml
a: ['dog', 'cat']
```

## Prepend to existing array
Given a sample.yml file of:
```yaml
a:
  - dog
```
then
```bash
yq '.a = ["cat"] + .a' sample.yml
```
will output
```yaml
a:
  - cat
  - dog
```

## Add new object to array
Given a sample.yml file of:
```yaml
a:
  - dog: woof
```
then
```bash
yq '.a + {"cat": "meow"}' sample.yml
```
will output
```yaml
- dog: woof
- cat: meow
```

## Relative append
Given a sample.yml file of:
```yaml
a:
  a1:
    b:
      - cat
  a2:
    b:
      - dog
  a3: {}
```
then
```bash
yq '.a[].b += ["mouse"]' sample.yml
```
will output
```yaml
a:
  a1:
    b:
      - cat
      - mouse
  a2:
    b:
      - dog
      - mouse
  a3:
    b:
      - mouse
```

## String concatenation
Given a sample.yml file of:
```yaml
a: cat
b: meow
```
then
```bash
yq '.a += .b' sample.yml
```
will output
```yaml
a: catmeow
b: meow
```

## Number addition - float
If the lhs or rhs are floats then the expression will be calculated with floats.

Given a sample.yml file of:
```yaml
a: 3
b: 4.9
```
then
```bash
yq '.a = .a + .b' sample.yml
```
will output
```yaml
a: 7.9
b: 4.9
```

## Number addition - int
If both the lhs and rhs are ints then the expression will be calculated with ints.

Given a sample.yml file of:
```yaml
a: 3
b: 4
```
then
```bash
yq '.a = .a + .b' sample.yml
```
will output
```yaml
a: 7
b: 4
```

## Increment numbers
Given a sample.yml file of:
```yaml
a: 3
b: 5
```
then
```bash
yq '.[] += 1' sample.yml
```
will output
```yaml
a: 4
b: 6
```

## Date addition
You can add durations to dates. Assumes RFC3339 date time format, see [date-time operators](https://mikefarah.gitbook.io/yq/operators/date-time-operators) for more information.

Given a sample.yml file of:
```yaml
a: 2021-01-01T00:00:00Z
```
then
```bash
yq '.a += "3h10m"' sample.yml
```
will output
```yaml
a: 2021-01-01T03:10:00Z
```

## Date addition - custom format
You can add durations to dates. See [date-time operators](https://mikefarah.gitbook.io/yq/operators/date-time-operators) for more information.

Given a sample.yml file of:
```yaml
a: Saturday, 15-Dec-01 at 2:59AM GMT
```
then
```bash
yq 'with_dtf("Monday, 02-Jan-06 at 3:04PM MST", .a += "3h1m")' sample.yml
```
will output
```yaml
a: Saturday, 15-Dec-01 at 6:00AM GMT
```

## Add to null
Adding to null simply returns the rhs

Running
```bash
yq --null-input 'null + "cat"'
```
will output
```yaml
cat
```

## Add maps to shallow merge
Adding objects together shallow merges them. Use `*` to deeply merge.

Given a sample.yml file of:
```yaml
a:
  thing:
    name: Astuff
    value: x
  a1: cool
b:
  thing:
    name: Bstuff
    legs: 3
  b1: neat
```
then
```bash
yq '.a += .b' sample.yml
```
will output
```yaml
a:
  thing:
    name: Bstuff
    legs: 3
  a1: cool
  b1: neat
b:
  thing:
    name: Bstuff
    legs: 3
  b1: neat
```

## Custom types: that are really strings
When custom tags are encountered, yq will try to decode the underlying type.

Given a sample.yml file of:
```yaml
a: !horse cat
b: !goat _meow
```
then
```bash
yq '.a += .b' sample.yml
```
will output
```yaml
a: !horse cat_meow
b: !goat _meow
```

## Custom types: that are really numbers
When custom tags are encountered, yq will try to decode the underlying type.

Given a sample.yml file of:
```yaml
a: !horse 1.2
b: !goat 2.3
```
then
```bash
yq '.a += .b' sample.yml
```
will output
```yaml
a: !horse 3.5
b: !goat 2.3
```

