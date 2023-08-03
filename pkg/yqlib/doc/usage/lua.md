
## Basic example
Given a sample.yml file of:
```yaml
hello: world
? look: non-string keys
: True
numbers: [123,456]

```
then
```bash
yq -o=lua '.' sample.yml
```
will output
```lua
return {["hello"]="world";[{["look"]="non-string keys";}]=true;["numbers"]={123,456,};};
```

