
## Minimum int
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
yq 'min' sample.yml
```
will output
```yaml
6
```

## Minimum string
Given a sample.yml file of:
```yaml
- foo
- bar
- baz
```
then
```bash
yq 'min' sample.yml
```
will output
```yaml
bar
```

## Minimum of empty
Given a sample.yml file of:
```yaml
[]
```
then
```bash
yq 'min' sample.yml
```
will output
```yaml
```

