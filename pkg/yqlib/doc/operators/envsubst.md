# Envsubst

This operator is used to replace environment variables in strings using [envsubst](https://github.com/a8m/envsubst).

To replace environment variables across all values in a document, this can be used with the recursive descent operator
as follows:

```bash
yq eval '(.. | select(tag == "!!str)) |= envsubst' file.yaml
```

## Replace strings with envsubst
Running
```bash
myenv="cat" yq eval --null-input '"the ${myenv} meows" | envsubst'
```
will output
```yaml
the cat meows
```

## Replace strings with envsubst, missing variables
Running
```bash
myenv="cat" yq eval --null-input '"the ${myenvnonexisting} meows" | envsubst'
```
will output
```yaml
the  meows
```

## Replace strings with envsubst, missing variables with defaults
Running
```bash
myenv="cat" yq eval --null-input '"the ${myenvnonexisting-dog} meows" | envsubst'
```
will output
```yaml
the dog meows
```

## Replace string environment variable in document
Given a sample.yml file of:
```yaml
v: ${myenv}
```
then
```bash
myenv="cat meow" yq eval '.v |= envsubst' sample.yml
```
will output
```yaml
v: cat meow
```

