
## Number modulo - int
If the lhs and rhs are ints then the expression will be calculated with ints.

Given a sample.yml file of:
```yaml
a: 13
b: 2
```
then
```bash
yq '.a = .a % .b' sample.yml
```
will output
```yaml
a: 1
b: 2
```

## Number modulo - float
If the lhs or rhs are floats then the expression will be calculated with floats.

Given a sample.yml file of:
```yaml
a: 12
b: 2.5
```
then
```bash
yq '.a = .a % .b' sample.yml
```
will output
```yaml
a: !!float 2
b: 2.5
```

