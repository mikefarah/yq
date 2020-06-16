---
description: Validate a given yaml file
---

# Validate

```text
yq v <yaml_file|->
```

Validates the given yaml file, does nothing if its valid, otherwise it will print errors to Stderr and exit with a non 0 exit code. This works like the [read command](read.md) - but does not take a path expression and does not print the yaml if it is valid.

## Basic - valid

```text
yq v valid.yaml
```

This will not print anything, and finish with a successful \(0\) exit code.

## Basic - invalid, from stdin

```text
echo '[1234' | yq v -
```

This will print the parsing error to stderr:

```bash
﻿﻿10:43:09 main [ERRO] yaml: line 1: did not find expected ',' or ']'
```

And return a error exit code \(1\)

## Multiple documents

Like other commands, by default the validate command will only run against the first document in the yaml file. Note that when running against other specific document indexes, _all previous documents will also be validated._

### Validating a single document

```bash
yq v -d1 multidoc.yml
```

This will validate both document 0 and document 1 \(but not document 2\)

### Validating all documents

```bash
yq v -d'*' multidoc.yml
```

This will validate all documents in the yaml file. Note that \* is quoted to avoid the CLI from processing the wildcard.

