# File Operators

File operators are most often used with merge when needing to merge specific files together. Note that when doing this, you will need to use `eval-all` to ensure all yaml documents are loaded into memory before performing the merge (as opposed to `eval` which runs the expression once per document).

Note that the `fileIndex` operator has a short alias of `fi`.

## Merging files
Note the use of eval-all to ensure all documents are loaded into memory.
```bash
yq eval-all 'select(fi == 0) * select(filename == "file2.yaml")' file1.yaml file2.yaml
```

{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Get filename
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq 'filename' sample.yml
```
will output
```yaml
sample.yml
```

## Get file index
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq 'file_index' sample.yml
```
will output
```yaml
0
```

## Get file indices of multiple documents
Given a sample.yml file of:
```yaml
a: cat
```
And another sample another.yml file of:
```yaml
a: cat
```
then
```bash
yq eval-all 'file_index' sample.yml another.yml
```
will output
```yaml
0
---
1
```

## Get file index alias
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq 'fi' sample.yml
```
will output
```yaml
0
```

