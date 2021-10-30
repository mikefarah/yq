---
description: >-
  Evaluates the given expression against each yaml document in each file, in
  sequence
---

# Evaluate

## Usage:&#x20;

```bash
yq eval [expression] [yaml_file1]... [flags]
```

Aliases: `eval, e`

Note that you can pass in `-` as a filename to pipe from STDIN.

## Examples:

```bash
# runs the expression against each file, in series
yq e '.a.b | length' f1.yml f2.yml 

# '-' will pipe from STDIN
cat file.yml | yq e '.a.b' f1.yml -  f2.yml

# prints out the file
yq e sample.yaml 
cat sample.yml | yq e

# prints a new yaml document
yq e -n '.a.b.c = "cat"' 

# updates file.yaml directly
yq e '.a.b = "cool"' -i file.yaml 
```



## Flags:

```
  -h, --help          help for eval
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
