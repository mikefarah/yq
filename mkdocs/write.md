```
yq w <yaml_file> <path> <new value>
```

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

### Splat
Given a sample.yaml file of:
```yaml
---
bob:
  item1:
    cats: bananas
  item2:
    cats: apples
  thing:
    cats: oranges
```
then
```bash
yq w sample.yaml bob.*.cats meow
```
will output:
```yaml
---
bob:
  item1:
    cats: meow
  item2:
    cats: meow
  thing:
    cats: meow
```

### Prefix Splat
Given a sample.yaml file of:
```yaml
---
bob:
  item1:
    cats: bananas
  item2:
    cats: apples
  thing:
    cats: oranges
```
then
```bash
yq w sample.yaml bob.item*.cats meow
```
will output:
```yaml
---
bob:
  item1:
    cats: meow
  item2:
    cats: meow
  thing:
    cats: oranges
```

### Array Splat
Given a sample.yaml file of:
```yaml
---
bob:
- cats: bananas
- cats: apples
- cats: oranges
```
then
```bash
yq w sample.yaml bob[*].cats meow
```
will output:
```yaml
---
bob:
- cats: meow
- cats: meow
- cats: meow
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

### Multiple Documents - update a single document
Given a sample.yaml file of:
```yaml
something: else
---
b:
  c: 2
```
then
```bash
yq w -d1 sample.yaml b.c 5
```
will output:
```yaml
something: else
---
b:
  c: 5
```

### Multiple Documents - update all documents
Given a sample.yaml file of:
```yaml
something: else
---
b:
  c: 2
```
then
```bash
yq w -d'*' sample.yaml b.c 5
```
will output:
```yaml
something: else
b:
  c: 5
---
b:
  c: 5
```

Note that '*' is in quotes to avoid being interpreted by your shell.

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

{!snippets/niche.md!}
