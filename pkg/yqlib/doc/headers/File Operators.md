The file operator is used to filter based on filename. This is most often used with merge when needing to merge specific files together.

```bash
yq eval 'filename == "file1.yaml" * fileIndex == 0' file1.yaml file2.yaml
```