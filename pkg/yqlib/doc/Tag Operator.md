The tag operator can be used to get or set the tag of ndoes (e.g. `!!str`, `!!int`, `!!bool`).
## Examples
### Get tag
Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
f: []
```
then
```bash
yq eval '.. | tag' sample.yml
```
will output
```yaml
!!map
!!str
!!int
!!float
!!bool
!!seq
```

### Convert numbers to strings
Given a sample.yml file of:
```yaml
a: cat
b: 5
c: 3.2
e: true
```
then
```bash
yq eval '(.. | select(tag == "!!int")) tag= "!!str"' sample.yml
```
will output
```yaml
a: cat
b: "5"
c: 3.2
e: true
```

