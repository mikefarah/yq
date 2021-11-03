# Env Variable Operators

This operator is used to handle environment variables usage in path expressions. While environment variables can, of course, be passed in via your CLI with string interpolation, this often comes with complex quote escaping and can be tricky to write and read. Note that there are two forms, `env` which will parse the environment variable as a yaml (be it a map, array, string, number of boolean) and `strenv` which will always parse the argument as a string.

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

