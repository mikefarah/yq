---
description: New features and breaking changes
---

# Upgrading from V2

## New Features

* Keeps yaml comments and formatting, can specify yaml [tags](usage/value-parsing.md#using-the-tag-field-to-override) when updating.
* Handles anchors!
* Can print out matching paths and values when splatting, more info [here](commands/read.md#printing-matching-paths).
* JSON output works for all commands! Yaml files with multiple documents are printed out as one JSON document per line, more info [here](usage/convert.md)
* Deep splat \(`**`\) to match arbitrary paths and match nodes by their children, more info [here](usage/path-expressions.md)

## Breaking Changes

### Parsing values from the CLI

In V3 users are able to better control how values are treated when updating YAML by using a new `--tag` argument \(see more info [here](usage/value-parsing.md)\). A result of this however, is that quoting values, e.g. "true" will no longer have an effect on how the value is interpreted like it did in V2.

For instance, to get the _string_ "true" into a yaml file:

V2:

```text
yq n a.path '"true"'
```

V3

```text
yq n a.path --tag '!!str' true
```

### Reading paths that don't exist

In V2 this would return null, V3 does not return anything.

Similarly, reading null yaml values `null`, `~` and , V2 returns null whereas V3 returns the values as is.

This is a result of taking effort not to format values coming in and out of the original YAML.



### Update scripts file format has changed to be more powerful.

Comments can be added, and delete commands have been introduced.

V2

```text
b.e[+].name: Mike Farah
```

V3

```yaml
- command: update 
  path: b.e[+].thing
  value:
    #great 
    things: frog # wow!
- command: delete
  path: b.d
```

### Reading and splatting, matching results are printed once per line.

e.g:

```yaml
parent:
  childA: 
    no: matches here
  childB:
    there: matches
    hi: no match
    there2: also matches
```

```text
yq r sample.yaml 'parent.*.there*'
```

V2

```text
- null
- - matches
  - also matches
```

V3

```text
matches
also matches
```

### Converting JSON to YAML

As JSON is a subset of YAML, and `yq` now preserves the formatting of the passed in document, you will most likely need to use the `--prettyPrint` flag to format the JSON document as idiomatic YAML. See [Working with JSON](usage/convert.md#json-to-yaml) for more info.





