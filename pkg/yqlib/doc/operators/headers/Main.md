# NAME
  *yq* is a portable command-line YAML processor

# SYNOPSIS 

yq [eval/eval-all] [expression] files..

eval/e  - (default) Apply the expression to each document in each yaml file in sequence

eval-all/ea - Loads all yaml documents of all yaml files and runs expression once

# DESCRIPTION

a lightweight and portable command-line YAML processor. `yq` uses [jq](https://github.com/stedolan/jq) like syntax but works with yaml files as well as json. It doesn't yet support everything `jq` does - but it does support the most common operations and functions, and more is being added continuously.

This documentation is also available at https://mikefarah.gitbook.io/yq/
# QUICK GUIDE 

## Read a value:
```bash
yq '.a.b[0].c' file.yaml
```

## Pipe from STDIN:
```bash
cat file.yaml | yq '.a.b[0].c'
```

## Update a yaml file, in place
```bash
yq -i '.a.b[0].c = "cool"' file.yaml
```

## Update using environment variables
```bash
NAME=mike yq -i '.a.b[0].c = strenv(NAME)' file.yaml
```

## Merge multiple files
```
yq ea '. as $item ireduce ({}; . * $item )' path/to/*.yml
```
Note the use of `ea` to evaluate all files at once (instead of in sequence.)

## Multiple updates to a yaml file
```bash
yq -i '
  .a.b[0].c = "cool" |
  .x.y.z = "foobar" |
  .person.name = strenv(NAME)
' file.yaml
```

See the [documentation](https://mikefarah.gitbook.io/yq/) for more.

# KNOWN ISSUES / MISSING FEATURES
- `yq` attempts to preserve comment positions and whitespace as much as possible, but it does not handle all scenarios (see https://github.com/go-yaml/yaml/tree/v3 for details)
- Powershell has its own...opinions: https://mikefarah.gitbook.io/yq/usage/tips-and-tricks#quotes-in-windows-powershell

# BUGS / ISSUES / FEATURE REQUESTS

Please visit the GitHub page https://github.com/mikefarah/yq/.

