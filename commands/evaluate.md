---
description: >-
  Evaluates the given expression against each yaml document in each file, in
  sequence
---

# Evaluate

Note that (as of 4.18.1) this is the default command when none is supplied to yq.

## Usage:&#x20;

```bash
yq eval [expression] [yaml_file1]... [flags]
```

Aliases: `eval, e`

Note that you can pass in `-` as a filename to pipe from STDIN.

## Examples:

```bash
# runs the expression against each file, in series
yq '.a.b | length' f1.yml f2.yml 

# '-' will pipe from STDIN
cat file.yml | yq '.a.b' f1.yml -  f2.yml

# prints out the file
yq sample.yaml 
cat sample.yml | yq e

# prints a new yaml document
yq -n '.a.b.c = "cat"' 

# updates file.yaml directly
yq '.a.b = "cool"' -i file.yaml 
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
