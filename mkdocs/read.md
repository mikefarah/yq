```
yaml r <yaml file> <path>
```

### Basic
Given a sample.yaml file of:
```yaml
b:
  c: 2
```
then
```bash
yaml r sample.yaml b.c
```
will output the value of '2'.

### From Stdin
Given a sample.yaml file of:
```bash
cat sample.yaml | yaml r - b.c
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
```
then
```bash
yaml r sample.yaml bob.*.cats
```
will output
```yaml
- bananas
- apples
```

### Handling '.' in the yaml key
Given a sample.yaml file of:
```yaml
b.x:
  c: 2
```
then
```bash
yaml r sample.yaml \"b.x\".c
```
will output the value of '2'.

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
yaml r sample.yaml b.e[1].name
```
will output 'sam'

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
yaml r sample.yaml b.e[*].name
```
will output:
```
- fred
- sam
```
