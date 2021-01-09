
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

