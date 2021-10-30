---
description: Flags to control yaml and json output format
---

# Output format

These flags are available for all `yq` commands.&#x20;

## Color

By default, `yq` prints with colours if it detects a terminal. You can manully this by using either

The `--colors/-C`flag to force print with colors.&#x20;

The \``--no-colors/-M` flag to force print without colours

## Pretty Print

To print out idiomatic `yaml` use the `--prettyPrint/-P` flag. Note that this is shorthand for using the [style](broken-reference) operator `... style=""`

## Indent

Use the indent flag `--indent/-I` to control the number of spaces used for indentation. This also works for JSON output. The default value is 2.&#x20;

Note that lists are indented at the same level as the map key at indent level 2, but are more deeply indented at indent level 4 and greater. This is (currently) a quirk of the underlying [yaml parser](https://github.com/go-yaml/yaml/tree/v3).

Given:

```
apples:
  collection:
  - name: Green
  - name: Blue
  favourite: Pink Lady
```

Then:

```
yq e -I4 sample.yaml
```

Will print out:

```yaml
apples:
    collection:
      - name: Green
      - name: Blue
    favourite: Pink Lady
```

This also works with json

```
yq e -j -I4 sample.yaml
```

yields

```javascript
{
    "apples": {
        "collection": [
            {
                "name": "Green"
            },
            {
                "name": "Blue"
            }
        ],
        "favourite": "Pink Lady"
    }
}
```

## Unwrap scalars

By default scalar values are 'unwrapped', that is only their value is printed (except when outputting as JSON). To print out the node as-is, with the original formatting an any comments pass in `--unwrapScalar=false`

Given data.yml:

```yaml
a: "Things" # cool stuff
```

Then:

`yq e --unwrapScalar=false '.a' data.yml`

Will yield:

```yaml
"Things" # cool stuff
```

where as without setting the flag to false you would get:

```yaml
Things
```

