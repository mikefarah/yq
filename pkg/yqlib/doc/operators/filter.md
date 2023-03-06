
## Filter array
Given a sample.yml file of:
```yaml
- 1
- 2
- 3
```
then
```bash
yq 'filter(. < 3)' sample.yml
```
will output
```yaml
- 1
- 2
```

