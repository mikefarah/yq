
## String split
Given a sample.yml file of:
```yaml
a: cat_meow
b: _
```
then
```bash
yq '.c = .a / .b' sample.yml
```
will output
```yaml
a: cat_meow
b: _
c:
  - cat
  - meow
```

## Number division
The result during division is calculated as a float

Given a sample.yml file of:
```yaml
a: 12
b: 2.5
```
then
```bash
yq '.a = .a / .b' sample.yml
```
will output
```yaml
a: 4.8
b: 2.5
```

## Number division by zero
Dividing by zero results in +Inf or -Inf

Given a sample.yml file of:
```yaml
a: 1
b: -1
```
then
```bash
yq '.a = .a / 0 | .b = .b / 0' sample.yml
```
will output
```yaml
a: !!float +Inf
b: !!float -Inf
```

