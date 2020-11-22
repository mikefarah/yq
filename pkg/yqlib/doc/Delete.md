Deletes matching entries in maps or arrays.
## Delete entry in map
Given a sample.yml file of:
```yaml
a: cat
b: dog
```
then
```bash
yq eval 'del(.b)' sample.yml
```
will output
```yaml
a: cat
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

## Delete no matches
Given a sample.yml file of:
```yaml
a: cat
b: dog
```
then
```bash
yq eval 'del(.c)' sample.yml
```
will output
```yaml
a: cat
b: dog
```

## Delete matching entries
Given a sample.yml file of:
```yaml
a: cat
b: dog
c: bat
```
then
```bash
yq eval 'del( .[] | select(. == "*at") )' sample.yml
```
will output
```yaml
b: dog
```

