
## Maximum int
Given a sample.yml file of:
```yaml
- 99
- 16
- 12
- 6
- 66
```
then
```bash
yq 'max' sample.yml
```
will output
```yaml
99
```

## Maximum string
Given a sample.yml file of:
```yaml
- foo
- bar
- baz
```
then
```bash
yq 'max' sample.yml
```
will output
```yaml
foo
```

## Maximum of empty
Given a sample.yml file of:
```yaml
[]
```
then
```bash
yq 'max' sample.yml
```
will output
```yaml
```

