# Add

Add behaves differently according to the type of the LHS:
* arrays: concatenate
* number scalars: arithmetic addition
* string scalars: concatenate

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

## Append to array
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

## Append another array using +=
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
yq eval '.a += .b' sample.yml
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
yq eval '.a[].b += ["mouse"]' sample.yml
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
  a3: {b: [mouse]}
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

## Increment numbers
Given a sample.yml file of:
```yaml
a: 3
b: 5
```
then
```bash
yq eval '.[] += 1' sample.yml
```
will output
```yaml
a: 4
b: 6
```

## Add to null
Adding to null simply returns the rhs

Running
```bash
yq eval --null-input 'null + "cat"'
```
will output
```yaml
cat
```

