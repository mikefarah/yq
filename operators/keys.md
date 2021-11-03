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

