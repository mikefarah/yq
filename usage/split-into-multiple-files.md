# Split into multiple files

`yq` can split out the results into multiple files with the `--split-exp/s` flag. You will need to give this flag an expression (that returns a string), this will be used as the filename for each result. In this expression, you can use `$index` to represent the result index in the name, if desired.

## Split documents into files

Given a file like

```yaml
a: test_doc1
--- 
a: test_doc2
```

Then running:

```bash
yq -s '.a' myfile.yml
```

will result in two files:

test\_doc1.yml:

```yaml
a: test_doc1
```

test\_doc2.yml:

```yaml
---
a: test_doc2
```

TIP: if you don't want the leading document separators (`---`), then run with the `--no-doc` flag.

## Split documents into files, using index

This is like the example above, but we'll use `$index` for the filename. Note that this variable is only defined for the `--split-exp/s` flag.

```
yq -s '"file_" + $index' myfile.yml
```

This will create two files, `file_0.yml` and `file_1.yml`.

## Split single document into files

You can also split results into separate files. Notice

```yaml
- name: bob
  age: 23
- name: tim
  age: 17
```

Then, by splatting the array into individual results, we can split the content into several files:

```bash
yq '.[]' file.yml -s '"user_" + .name'
```

will result in two files:

user\_bob.yml:

```yaml
name: bob
age: 23
```

user\_tim.yml:

```yaml
name: tim
age: 17
```

