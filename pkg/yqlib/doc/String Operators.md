
## Join strings
Given a sample.yml file of:
```yaml
- cat
- meow
- 1
- null
- true
```
then
```bash
yq eval 'join("; ")' sample.yml
```
will output
```yaml
cat; meow; 1; ; true
```

