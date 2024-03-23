
## Omit keys from map
Note that the order of the keys matches the omit order and non existent keys are skipped.

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
yq '.myMap |= omit(["hamster", "cat", "goat"])' sample.yml
```
will output
```yaml
myMap:
  dog: bark
  thing: hamster
```

## Omit indices from array
Note that the order of the indices matches the omit order and non existent indices are skipped.

Given a sample.yml file of:
```yaml
- cat
- leopard
- lion
```
then
```bash
yq 'omit([2, 0, 734, -5])' sample.yml
```
will output
```yaml
- leopard
```

