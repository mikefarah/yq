
## Number subtraction - float
If the lhs or rhs are floats then the expression will be calculated with floats.

Given a sample.yml file of:
```yaml
a: 3
b: 4.5
```
then
```bash
yq eval '.a = .a - .b' sample.yml
```
will output
```yaml
a: -1.5
b: 4.5
```

## Number subtraction - float
If the lhs or rhs are floats then the expression will be calculated with floats.

Given a sample.yml file of:
```yaml
a: 3
b: 4.5
```
then
```bash
yq eval '.a = .a - .b' sample.yml
```
will output
```yaml
a: -1.5
b: 4.5
```

## Number subtraction - int
If both the lhs and rhs are ints then the expression will be calculated with ints.

Given a sample.yml file of:
```yaml
a: 3
b: 4
```
then
```bash
yq eval '.a = .a - .b' sample.yml
```
will output
```yaml
a: -1
b: 4
```

## Decrement numbers
Given a sample.yml file of:
```yaml
a: 3
b: 5
```
then
```bash
yq eval '.[] -= 1' sample.yml
```
will output
```yaml
a: 2
b: 4
```

