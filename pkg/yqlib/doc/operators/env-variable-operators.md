# Env Variable Operators

These operators are used to handle environment variables usage in expressions and documents. While environment variables can, of course, be passed in via your CLI with string interpolation, this often comes with complex quote escaping and can be tricky to write and read. 

There are three operators:
-  `env` which takes a single environment variable name and parse the variable as a yaml node (be it a map, array, string, number of boolean) 
- `strenv` which also takes a single environment variable name, and always parses the variable as a string.
- `envsubst` which you pipe strings into and it interpolates environment variables in strings using [envsubst](https://github.com/a8m/envsubst). 


## Tip
To replace environment variables across all values in a document, `envsubst` can be used with the recursive descent operator
as follows:

```bash
yq eval '(.. | select(tag == "!!str")) |= envsubst' file.yaml
```


## Read string environment variable
Running
```bash
myenv="cat meow" yq eval --null-input '.a = env(myenv)'
```
will output
```yaml
a: cat meow
```

## Read boolean environment variable
Running
```bash
myenv="true" yq eval --null-input '.a = env(myenv)'
```
will output
```yaml
a: true
```

## Read numeric environment variable
Running
```bash
myenv="12" yq eval --null-input '.a = env(myenv)'
```
will output
```yaml
a: 12
```

## Read yaml environment variable
Running
```bash
myenv="{b: fish}" yq eval --null-input '.a = env(myenv)'
```
will output
```yaml
a: {b: fish}
```

## Read boolean environment variable as a string
Running
```bash
myenv="true" yq eval --null-input '.a = strenv(myenv)'
```
will output
```yaml
a: "true"
```

## Read numeric environment variable as a string
Running
```bash
myenv="12" yq eval --null-input '.a = strenv(myenv)'
```
will output
```yaml
a: "12"
```

## Dynamic key lookup with environment variable
Given a sample.yml file of:
```yaml
cat: meow
dog: woof
```
then
```bash
myenv="cat" yq eval '.[env(myenv)]' sample.yml
```
will output
```yaml
meow
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

