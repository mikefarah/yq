# Keys

Use the `keys` operator to return map keys or array indices. 

## Map keys
Given a sample.yml file of:
```yaml
dog: woof
cat: meow
```
then
```bash
yq eval 'keys' sample.yml
```
will output
```yaml
- dog
- cat
```

## Array keys
Given a sample.yml file of:
```yaml
- apple
- banana
```
then
```bash
yq eval 'keys' sample.yml
```
will output
```yaml
- 0
- 1
```

## Retrieve array key
Given a sample.yml file of:
```yaml
- 1
- 2
- 3
```
then
```bash
yq eval '.[1] | key' sample.yml
```
will output
```yaml
1
```

## Retrieve map key
Given a sample.yml file of:
```yaml
a: thing
```
then
```bash
yq eval '.a | key' sample.yml
```
will output
```yaml
a
```

## No key
Given a sample.yml file of:
```yaml
{}
```
then
```bash
yq eval 'key' sample.yml
```
will output
```yaml
```

## Update map key
Given a sample.yml file of:
```yaml
a:
  x: 3
  y: 4
```
then
```bash
yq eval '(.a.x | key) = "meow"' sample.yml
```
will output
```yaml
a:
  meow: 3
  y: 4
```

## Get comment from map key
Given a sample.yml file of:
```yaml
a:
  # comment on key
  x: 3
  y: 4
```
then
```bash
yq eval '.a.x | key | headComment' sample.yml
```
will output
```yaml
comment on key
```

