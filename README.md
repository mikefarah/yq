# yq

a lightweight and portable command-line YAML processor. `yq` uses [jq](https://github.com/stedolan/jq) like syntax but works with yaml files as well as json. It doesn't yet support everything `jq` does - but it does support the most common operations and functions, and more is being added continuously.

yq is written in go - so you can download a dependency free binary for your platform and you are good to go! If you prefer there are a variety of package managers that can be used as well as docker, all listed below.

## Quick Usage Guide

Read a value:
```bash
yq '.a.b[0].c' file.yaml
```

Pipe from STDIN:
```bash
cat file.yaml | yq '.a.b[0].c'
```

Update a yaml file, inplace
```bash
yq -i '.a.b[0].c = "cool"' file.yaml
```

Update using environment variables
```bash
NAME=mike yq -i '.a.b[0].c = strenv(NAME)' file.yaml
```

Merge multiple files
```
yq ea '. as $item ireduce ({}; . * $item )' path/to/*.yml
```

Multiple updates to a yaml file
```bash
yq -i '
  .a.b[0].c = "cool" |
  .x.y.z = "foobar" |
  .person.name = strenv(NAME)
' file.yaml
```

Take a look at the discussions for [common questions](https://github.com/mikefarah/yq/discussions/categories/q-a), and [cool ideas](https://github.com/mikefarah/yq/discussions/categories/show-and-tell)

## Install

See the [github page](https://github.com/mikefarah/yq/#install) for the various ways you can install and use `yq`

## Known Issues / Missing Features
- `yq` attempts to preserve comment positions and whitespace as much as possible, but it does not handle all scenarios (see https://github.com/go-yaml/yaml/tree/v3 for details)
- Powershell has its own...[opinions on quoting yq](https://mikefarah.gitbook.io/yq/usage/tips-and-tricks#quotes-in-windows-powershell)

See [tips and tricks](https://mikefarah.gitbook.io/yq/usage/tips-and-tricks) for more common problems and solutions.