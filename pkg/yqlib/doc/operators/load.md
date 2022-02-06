# Load

The `load`/`strload` operator allows you to load in content from another file referenced in your yaml document.

Note that you can use string operators like `+` and `sub` to modify the value in the yaml file to a path that exists in your system.

Use `strload` to load text based content as a string block, and `load` to interpret the file as yaml.

Lets say there is a file `../../examples/thing.yml`:

```yaml
a: apple is included
b: cool
```

{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Simple example
Given a sample.yml file of:
```yaml
myFile: ../../examples/thing.yml
```
then
```bash
yq 'load(.myFile)' sample.yml
```
will output
```yaml
a: apple is included
b: cool.
```

## Replace node with referenced file
Note that you can modify the filename in the load operator if needed.

Given a sample.yml file of:
```yaml
something:
  file: thing.yml
```
then
```bash
yq '.something |= load("../../examples/" + .file)' sample.yml
```
will output
```yaml
something:
  a: apple is included
  b: cool.
```

## Replace _all_ nodes with referenced file
Recursively match all the nodes (`..`) and then filter the ones that have a 'file' attribute. 

Given a sample.yml file of:
```yaml
something:
  file: thing.yml
over:
  here:
    - file: thing.yml
```
then
```bash
yq '(.. | select(has("file"))) |= load("../../examples/" + .file)' sample.yml
```
will output
```yaml
something:
  a: apple is included
  b: cool.
over:
  here:
    - a: apple is included
      b: cool.
```

## Replace node with referenced file as string
This will work for any text based file

Given a sample.yml file of:
```yaml
something:
  file: thing.yml
```
then
```bash
yq '.something |= strload("../../examples/" + .file)' sample.yml
```
will output
```yaml
something: |-
  a: apple is included
  b: cool.
```

