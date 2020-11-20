The file operator is used to filter based on filename. This is most often used with merge when needing to merge specific files together.

```bash
yq eval 'filename == "file1.yaml" * fileIndex == 0' file1.yaml file2.yaml
```
## Examples
### Get filename
Given a sample.yml file of:
```yaml
'': null
```
then
```bash
yq eval 'filename' sample.yml
```
will output
```yaml
sample.yaml
```

### Get file index
Given a sample.yml file of:
```yaml
'': null
```
then
```bash
yq eval 'fileIndex' sample.yml
```
will output
```yaml
73
```

