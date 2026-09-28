
## Check key exists in map using variable binding
Given a sample.yml file of:
```yaml
a: 1
b: 2
c: 3
```
then
```bash
yq '. as $m | "a" | in($m)' sample.yml
```
will output
```yaml
true
```

## Check key does not exist in map
Given a sample.yml file of:
```yaml
a: 1
b: 2
c: 3
```
then
```bash
yq '. as $m | "d" | in($m)' sample.yml
```
will output
```yaml
false
```

## Check value exists in array
Given a sample.yml file of:
```yaml
- Tool
- Food
- Flower
```
then
```bash
yq '. as $m | "Food" | in($m)' sample.yml
```
will output
```yaml
true
```

## Check value does not exist in array
Given a sample.yml file of:
```yaml
- Tool
- Food
- Flower
```
then
```bash
yq '. as $m | "Animal" | in($m)' sample.yml
```
will output
```yaml
false
```

## Check in with select on array elements
Filter items whose type is in the given list

Given a sample.yml file of:
```yaml
- item: Pizza
  type: Food
- item: Rose
  type: Flower
- item: Hammer
  type: Tool
```
then
```bash
yq '.[] | select(.type | in(["Tool", "Food"]))' sample.yml
```
will output
```yaml
item: Pizza
type: Food
item: Hammer
type: Tool
```

## In with variable binding - found
Given a sample.yml file of:
```yaml
a: 1
b: 2
c: 3
```
then
```bash
yq '. as $m | "b" | in($m)' sample.yml
```
will output
```yaml
true
```

## In with variable binding - not found
Given a sample.yml file of:
```yaml
a: 1
b: 2
c: 3
```
then
```bash
yq '. as $m | "z" | in($m)' sample.yml
```
will output
```yaml
false
```

