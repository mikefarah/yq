
## Sort by string field
Given a sample.yml file of:
```yaml
- a: banana
- a: cat
- a: apple
```
then
```bash
yq eval 'sort_by(.a)' sample.yml
```
will output
```yaml
- a: apple
- a: banana
- a: cat
```

