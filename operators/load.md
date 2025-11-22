# Load

The load operators allows you to load in content from another file.

Note that you can use string operators like `+` and `sub` to modify the value in the yaml file to a path that exists in your system.

You can load files of the following supported types:

|Format | Load Operator |
| --- | --- |
| Yaml | load |
| XML | load_xml |
| Properties | load_props |
| Plain String | load_str |
| Base64 | load_base64 |

Note that load_base64 only works for base64 encoded utf-8 strings.

## Samples files for tests:

### yaml

`../../examples/thing.yml`:

```yaml
a: apple is included
b: cool
```

### xml
`small.xml`:

```xml
<this>is some xml</this>
```

### properties
`small.properties`:

```properties
this.is = a properties file
```

### base64
`base64.txt`:
```
bXkgc2VjcmV0IGNoaWxsaSByZWNpcGUgaXMuLi4u
```

## Disabling file operators
If required, you can use the `--security-disable-file-ops` to disable file operations.


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
yq '.something |= load_str("../../examples/" + .file)' sample.yml
```
will output
```yaml
something: |-
  a: apple is included
  b: cool.
```

## Load from XML
Given a sample.yml file of:
```yaml
cool: things
```
then
```bash
yq '.more_stuff = load_xml("../../examples/small.xml")' sample.yml
```
will output
```yaml
cool: things
more_stuff:
  this: is some xml
```

## Load from Properties
Given a sample.yml file of:
```yaml
cool: things
```
then
```bash
yq '.more_stuff = load_props("../../examples/small.properties")' sample.yml
```
will output
```yaml
cool: things
more_stuff:
  this:
    is: a properties file
```

## Merge from properties
This can be used as a convenient way to update a yaml document

Given a sample.yml file of:
```yaml
this:
  is: from yaml
  cool: ay
```
then
```bash
yq '. *= load_props("../../examples/small.properties")' sample.yml
```
will output
```yaml
this:
  is: a properties file
  cool: ay
```

## Load from base64 encoded file
Given a sample.yml file of:
```yaml
cool: things
```
then
```bash
yq '.more_stuff = load_base64("../../examples/base64.txt")' sample.yml
```
will output
```yaml
cool: things
more_stuff: my secret chilli recipe is....
```

## load() operation fails when security is enabled
Use `--security-disable-file-ops` to disable file operations for security.

Running
```bash
yq --null-input 'load("../../examples/thing.yml")'
```
will output
```bash
Error: file operations have been disabled
```

## load_str() operation fails when security is enabled
Use `--security-disable-file-ops` to disable file operations for security.

Running
```bash
yq --null-input 'load_str("../../examples/thing.yml")'
```
will output
```bash
Error: file operations have been disabled
```

## load_xml() operation fails when security is enabled
Use `--security-disable-file-ops` to disable file operations for security.

Running
```bash
yq --null-input 'load_xml("../../examples/small.xml")'
```
will output
```bash
Error: file operations have been disabled
```

## load_props() operation fails when security is enabled
Use `--security-disable-file-ops` to disable file operations for security.

Running
```bash
yq --null-input 'load_props("../../examples/small.properties")'
```
will output
```bash
Error: file operations have been disabled
```

## load_base64() operation fails when security is enabled
Use `--security-disable-file-ops` to disable file operations for security.

Running
```bash
yq --null-input 'load_base64("../../examples/base64.txt")'
```
will output
```bash
Error: file operations have been disabled
```

