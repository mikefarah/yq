Add behaves differently according to the type of the LHS:
- arrays: concatenate
- number scalars: arithmetic addition
- string scalars: concatenate

Use `+=` as append assign for things like increment. Note that `.a += .x` is equivalent to running `.a = .a + .x`.

## Concatenate and assign arrays
Given a sample.yml file of:
```yaml
a:
  val: thing
  b:
    - cat
    - dog
```
then
```bash
yq eval '.a.b += ["cow"]' sample.yml
```
will output
```yaml
a:
  val: thing
  b:
    - cat
    - dog
    - cow
```

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
yq eval '.a + .b' sample.yml
```
will output
```yaml
- 1
- 2
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
yq eval '.a + null' sample.yml
```
will output
```yaml
- 1
- 2
```

## Add new object to array
Given a sample.yml file of:
```yaml
a:
  - dog: woof
```
then
```bash
yq eval '.a + {"cat": "meow"}' sample.yml
```
will output
```yaml
- dog: woof
- cat: meow
```

## Add string to array
Given a sample.yml file of:
```yaml
a:
  - 1
  - 2
```
then
```bash
yq eval '.a + "hello"' sample.yml
```
will output
```yaml
- 1
- 2
- hello
```

## Update array (append)
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
yq eval '.a = .a + .b' sample.yml
```
will output
```yaml
a:
  - 1
  - 2
  - 3
  - 4
b:
  - 3
  - 4
```

## String concatenation
Given a sample.yml file of:
```yaml
a: cat
b: meow
```
then
```bash
yq eval '.a = .a + .b' sample.yml
```
will output
```yaml
a: catmeow
b: meow
```

## Relative string concatenation
Given a sample.yml file of:
```yaml
a: cat
b: meow
```
then
```bash
yq eval '.a += .b' sample.yml
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
yq eval '.a = .a + .b' sample.yml
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
yq eval '.a = .a + .b' sample.yml
```
will output
```yaml
a: 7
b: 4
```

## Increment number
Given a sample.yml file of:
```yaml
a: 3
```
then
```bash
yq eval '.a += 1' sample.yml
```
will output
```yaml
a: 4
```

