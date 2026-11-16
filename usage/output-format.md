---
description: Flags to control yaml and json output format
---

# Output format

These flags are available for all `yq` commands.

## Color

By default, `yq` prints with colours if it detects a terminal. You can manually change this by using either

The `--colors/-C` flag to force print with colors.

The `--no-colors/-M` flag to force print without colours

## Pretty Print

To print out idiomatic `yaml` use the `--prettyPrint/-P` flag. Note that this is shorthand for using the [style](https://mikefarah.gitbook.io/yq/operators/style) operator `... style=""`

## Indent

Use the indent flag `--indent/-I` to control the number of spaces used for indentation. This also works for JSON output. The default value is 2.

YAML sequences under a mapping key are indented by default. Use `--yaml-compact-seq-indent/-c` to count the `- ` marker toward the indentation. With compact indentation, a sequence starts at the same indentation as its key at indent level 2. This option maps to [Go YAML v4’s compact sequence indentation](https://pkg.go.dev/go.yaml.in/yaml/v4#WithCompactSeqIndent).

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
yq -I4 sample.yaml
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
yq -j -I4 sample.yaml
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

`yq --unwrapScalar=false '.a' data.yml`

Will yield:

```yaml
"Things" # cool stuff
```

where as without setting the flag to false you would get:

```yaml
Things
```

