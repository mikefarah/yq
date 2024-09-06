# Pick

Filter a map by the specified list of keys. Map is returned with the key in the order of the pick list.

Similarly, filter an array by the specified list of indices.

## Pick keys from map
Note that the order of the keys matches the pick order and non existent keys are skipped.

Given a sample.yml file of:
```yaml
myMap:
  cat: meow
  dog: bark
  thing: hamster
  hamster: squeak
```
then
```bash
yq '.myMap |= pick(["hamster", "cat", "goat"])' sample.yml
```
will output
```yaml
myMap:
  hamster: squeak
  cat: meow
```

## Pick keys from map, included all the keys
We create a map of the picked keys plus all the current keys, and run that through unique

Given a sample.yml file of:
```yaml
myMap:
  cat: meow
  dog: bark
  thing: hamster
  hamster: squeak
```
then
```bash
yq '.myMap |= pick( (["thing"] + keys) | unique)' sample.yml
```
will output
```yaml
myMap:
  thing: hamster
  cat: meow
  dog: bark
  hamster: squeak
```

## Pick indices from array
Note that the order of the indices matches the pick order and non existent indices are skipped.

Given a sample.yml file of:
```yaml
- cat
- leopard
- lion
```
then
```bash
yq 'pick([2, 0, 734, -5])' sample.yml
```
will output
```yaml
- lion
- cat
```

