```
yq r <yaml_file|json_file> <path>
```

{!snippets/works_with_json.md!}

### Basic
Given a sample.yaml file of:
```yaml
b:
  c: 2
```
then
```bash
yq r sample.yaml b.c
```
will output the value of '2'.

### From Stdin
Given a sample.yaml file of:
```bash
cat sample.yaml | yq r - b.c
```
will output the value of '2'.

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
yq r sample.yaml bob.*.cats
```
will output
```yaml
- bananas
- apples
- oranges
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
yq r sample.yaml bob.item*.cats
```
will output
```yaml
- bananas
- apples
```

### Multiple Documents - specify a single document
Given a sample.yaml file of:
```yaml
something: else
---
b:
  c: 2
```
then
```bash
yq r -d1 sample.yaml b.c
```
will output the value of '2'.

### Multiple Documents - read all documents
Reading all documents will return the result as an array. This can be converted to json using the '-j' flag if desired.

Given a sample.yaml file of:
```yaml
name: Fred
age: 22
---
name: Stella
age: 23
---
name: Android
age: 232
```
then
```bash
yq r -d'*' sample.yaml name
```
will output:
```
- Fred
- Stella
- Android
```

### Arrays
You can give an index to access a specific element:
e.g.: given a sample file of
```yaml
b:
  e:
    - name: fred
      value: 3
    - name: sam
      value: 4
```
then
```
yq r sample.yaml 'b.e[1].name'
```
will output 'sam'

Note that the path is in quotes to avoid the square brackets being interpreted by your shell.

### Array Splat
e.g.: given a sample file of
```yaml
b:
  e:
    - name: fred
      value: 3
    - name: sam
      value: 4
```
then
```
yq r sample.yaml 'b.e[*].name'
```
will output:
```
- fred
- sam
```
Note that the path is in quotes to avoid the square brackets being interpreted by your shell.

{!snippets/niche.md!}
