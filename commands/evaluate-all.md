---
description: >-
  Read all documents of all given yaml files into memory, then run the given
  expression once against the lot.
---

# Evaluate All

Evaluate All is most useful when needing to run expressions that depend on multiple yaml documents or files. Merge is probably the most common reason why evaluate all would be used. Note that `eval-all` consumes more memory than `evaluate`.

Like `evaluate` you can use `-` to pipe from STDIN.

## Usage

```bash
yq eval-all [expression] [yaml_file1]... [flags]
```

Aliases: `eval-all, ea`

## Examples

```bash
# merges f2.yml into f1.yml (inplace)
yq eval-all --inplace 'select(fileIndex == 0) * select(fileIndex == 1)' f1.yml f2.yml

# you can merge into a file, piping from STDIN
cat somefile.yml | yq eval-all --inplace 'select(fileIndex == 0) * select(fileIndex == 1)' f1.yml -
```

## Flags

```bash
  -h, --help          help for eval-all
  -C, --colors        force print with colors
  -e, --exit-status   set exit status if there are no matches or null or false is returned
  -I, --indent int    sets indent level for output (default 2)
  -i, --inplace       update the yaml file inplace of first yaml file given.
  -M, --no-colors     force print with no colors
  -N, --no-doc        Don't print document separators (---)
  -n, --null-input    Don't read input, simply evaluate the expression given. Useful for creating yaml docs from scratch.
  -j, --tojson        output as json. Set indent to 0 to print json in one line.
  -v, --verbose       verbose mode

```
