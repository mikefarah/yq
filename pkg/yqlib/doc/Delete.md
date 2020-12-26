Deletes matching entries in maps or arrays.
## Delete entry in map
Given a sample.yml file of:
```yaml
a: cat
b: dog
'': null
```
then
```bash
yq eval 'del(.b)' sample.yml
```
will output
```yaml
a: cat
'': null
```

## Delete nested entry in map
Given a sample.yml file of:
```yaml
a: {a1: fred, a2: frood}
'': null
```
then
```bash
yq eval 'del(.a.a1)' sample.yml
```
will output
```yaml
a: {a2: frood}
'': null
```

## Delete entry in array
Given a sample.yml file of:
```yaml
- 1
- 2
- 3
```
then
```bash
yq eval 'del(.[1])' sample.yml
```
will output
```yaml
- 1
- 3
```

## Delete nested entry in array
Given a sample.yml file of:
```yaml
- a: cat
  b: dog
  '': null
```
then
```bash
yq eval 'del(.[0].a)' sample.yml
```
will output
```yaml
- b: dog
  '': null
```

## Delete no matches
Given a sample.yml file of:
```yaml
a: cat
b: dog
'': null
```
then
```bash
yq eval 'del(.c)' sample.yml
```
will output
```yaml
a: cat
b: dog
'': null
```

## Delete matching entries
Given a sample.yml file of:
```yaml
a: cat
b: dog
c: bat
'': null
```
then
```bash
yq eval 'del( .[] | select(. == "*at") )' sample.yml
```
will output
```yaml
b: dog
'': null
```

