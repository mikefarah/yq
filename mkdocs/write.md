```
yq w <yaml_file|json_file> <path> <new value>
```
{!snippets/works_with_json.md!}

### To Stdout
Given a sample.yaml file of:
```yaml
b:
  c: 2
```
then
```bash
yq w sample.yaml b.c cat
```
will output:
```yaml
b:
  c: cat
```

### From STDIN
```bash
cat sample.yaml | yq w - b.c blah
```

### Adding new fields
Any missing fields in the path will be created on the fly.

Given a sample.yaml file of:
```yaml
b:
  c: 2
```
then
```bash
yq w sample.yaml b.d[0] "new thing"
```
will output:
```yaml
b:
  c: cat
  d:
    - new thing
```

### Appending value to an array field
Given a sample.yaml file of:
```yaml
b:
  c: 2
  d:
    - new thing
    - foo thing
```
then
```bash
yq w sample.yaml "b.d[+]" "bar thing"
```
will output:
```yaml
b:
  c: cat
  d:
    - new thing
    - foo thing
    - bar thing
```

Note that the path is in quotes to avoid the square brackets being interpreted by your shell.

### Updating files in-place
Given a sample.yaml file of:
```yaml
b:
  c: 2
```
then
```bash
yq w -i sample.yaml b.c cat
```
will update the sample.yaml file so that the value of 'c' is cat.


### Updating multiple values with a script
Given a sample.yaml file of:
```yaml
b:
  c: 2
  e:
    - name: Billy Bob
```
and a script update_instructions.yaml of:
```yaml
b.c: 3
b.e[0].name: Howdy Partner
```
then

```bash
yq w -s update_instructions.yaml sample.yaml
```
will output:
```yaml
b:
  c: 3
  e:
    - name: Howdy Partner
```

And, of course, you can pipe the instructions in using '-':
```bash
cat update_instructions.yaml | yq w -s - sample.yaml
```

### Values starting with a hyphen (or dash)
The flag terminator needs to be used to stop the app from attempting to parse the subsequent arguments as flags:

```
yq w -- my.path -3
```

will output
```yaml
my:
  path: -3
```

{!snippets/keys_with_dots.md!}
