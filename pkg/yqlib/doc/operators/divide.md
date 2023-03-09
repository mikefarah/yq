
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
The result during divison is calculated as a float

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

