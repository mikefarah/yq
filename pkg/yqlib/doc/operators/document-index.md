# Document Index

Use the `documentIndex` operator (or the `di` shorthand) to select nodes of a particular document.

{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Retrieve a document index
Given a sample.yml file of:
```yaml
a: cat
---
a: frog
```
then
```bash
yq '.a | document_index' sample.yml
```
will output
```yaml
0
---
1
```

## Retrieve a document index, shorthand
Given a sample.yml file of:
```yaml
a: cat
---
a: frog
```
then
```bash
yq '.a | di' sample.yml
```
will output
```yaml
0
---
1
```

## Filter by document index
Given a sample.yml file of:
```yaml
a: cat
---
a: frog
```
then
```bash
yq 'select(document_index == 1)' sample.yml
```
will output
```yaml
a: frog
```

## Filter by document index shorthand
Given a sample.yml file of:
```yaml
a: cat
---
a: frog
```
then
```bash
yq 'select(di == 1)' sample.yml
```
will output
```yaml
a: frog
```

## Print Document Index with matches
Given a sample.yml file of:
```yaml
a: cat
---
a: frog
```
then
```bash
yq '.a | ({"match": ., "doc": document_index})' sample.yml
```
will output
```yaml
match: cat
doc: 0
match: frog
doc: 1
```

