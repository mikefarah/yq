
## First matching element from array
Given a sample.yml file of:
```yaml
- a: banana
- a: cat
- a: apple
```
then
```bash
yq 'first(.a == "cat")' sample.yml
```
will output
```yaml
a: cat
```

## First matching element from array with multiple matches
Given a sample.yml file of:
```yaml
- a: banana
- a: cat
- a: apple
- a: cat
```
then
```bash
yq 'first(.a == "cat")' sample.yml
```
will output
```yaml
a: cat
```

## First matching element from array with numeric condition
Given a sample.yml file of:
```yaml
- a: 10
- a: 100
- a: 1
```
then
```bash
yq 'first(.a > 50)' sample.yml
```
will output
```yaml
a: 100
```

## First matching element from array with boolean condition
Given a sample.yml file of:
```yaml
- a: false
- a: true
- a: false
```
then
```bash
yq 'first(.a == true)' sample.yml
```
will output
```yaml
a: true
```

## First matching element from array with null values
Given a sample.yml file of:
```yaml
- a: null
- a: cat
- a: apple
```
then
```bash
yq 'first(.a != null)' sample.yml
```
will output
```yaml
a: cat
```

## First matching element from array with complex condition
Given a sample.yml file of:
```yaml
- a: dog
  b: 5
- a: cat
  b: 3
- a: apple
  b: 7
```
then
```bash
yq 'first(.b > 4)' sample.yml
```
will output
```yaml
a: dog
b: 5
```

## First matching element from map
Given a sample.yml file of:
```yaml
x:
  a: banana
y:
  a: cat
z:
  a: apple
```
then
```bash
yq 'first(.a == "cat")' sample.yml
```
will output
```yaml
a: cat
```

## First matching element from map with numeric condition
Given a sample.yml file of:
```yaml
x:
  a: 10
y:
  a: 100
z:
  a: 1
```
then
```bash
yq 'first(.a > 50)' sample.yml
```
will output
```yaml
a: 100
```

## First matching element from nested structure
Given a sample.yml file of:
```yaml
items:
  - a: banana
  - a: cat
  - a: apple
```
then
```bash
yq '.items | first(.a == "cat")' sample.yml
```
will output
```yaml
a: cat
```

## First matching element with no matches
Given a sample.yml file of:
```yaml
- a: banana
- a: cat
- a: apple
```
then
```bash
yq 'first(.a == "dog")' sample.yml
```
will output
```yaml
```

## First matching element from empty array
Given a sample.yml file of:
```yaml
[]
```
then
```bash
yq 'first(.a == "cat")' sample.yml
```
will output
```yaml
```

## First matching element from scalar node
Given a sample.yml file of:
```yaml
hello
```
then
```bash
yq 'first(. == "hello")' sample.yml
```
will output
```yaml
```

## First matching element from null node
Given a sample.yml file of:
```yaml
null
```
then
```bash
yq 'first(. == "hello")' sample.yml
```
will output
```yaml
```

## First matching element with string condition
Given a sample.yml file of:
```yaml
- a: banana
- a: cat
- a: apple
```
then
```bash
yq 'first(.a | test("^c"))' sample.yml
```
will output
```yaml
a: cat
```

## First matching element with length condition
Given a sample.yml file of:
```yaml
- a: hi
- a: hello
- a: world
```
then
```bash
yq 'first(.a | length > 4)' sample.yml
```
will output
```yaml
a: hello
```

## First matching element from array of strings
Given a sample.yml file of:
```yaml
- banana
- cat
- apple
```
then
```bash
yq 'first(. == "cat")' sample.yml
```
will output
```yaml
cat
```

## First matching element from array of numbers
Given a sample.yml file of:
```yaml
- 10
- 100
- 1
```
then
```bash
yq 'first(. > 50)' sample.yml
```
will output
```yaml
100
```

## First element with no filter from array
Given a sample.yml file of:
```yaml
- 10
- 100
- 1
```
then
```bash
yq 'first' sample.yml
```
will output
```yaml
10
```

## First element with no filter from array of maps
Given a sample.yml file of:
```yaml
- a: 10
- a: 100
```
then
```bash
yq 'first' sample.yml
```
will output
```yaml
a: 10
```

