---
description: Flags to control yaml and json output format
---

# Output format

These flags are available for all `yq` commands. 

## Colorize Output

Use the `--colors/-C`flag to print out yaml with colors. This does not work when outputing in JSON format.

## Pretty Print

Use the `--prettyPrint/-P` flag to enforce a formatting style for yaml documents. This is particularly useful when reading a json file \(which is a subset of yaml\) and wanting to format it in a more conventional yaml format.

Given:

```text
{
  "apples": [
    {
      "are": "great"
    }
  ]
}
```

Then:

```text
yq r --prettyPrint sample.json
```

Will print out:

```text
apples:
- are: great
```

This works in the same manner for yaml files:

```text
"apples": [are: great]
```

will format to:

```text
apples:
- are: great
```

## Indent

Use the indent flag `--indent/-I` to control the number of spaces used for indentation. This also works for JSON output. The default value is 2. 

Note that lists are indented at the same level as the map key at indent level 2, but are more deeply indented at indent level 4 and greater. This is \(currently\) a quirk of the underlying [yaml parser](https://github.com/go-yaml/yaml/tree/v3).

Given:

```text
apples:
  collection:
  - name: Green
  - name: Blue
  favourite: Pink Lady
```

Then:

```text
yq r -I4 sample.yaml
```

Will print out:

```text
apples:
    collection:
      - name: Green
      - name: Blue
    favourite: Pink Lady
```

With json, you must also specify the `--prettyPrint/-P` flag

```text
yq r -j -P -I4 sample.yaml
```

yields

```text
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

By default scalar values are 'unwrapped', that is only their value is printed \(except when outputting as JSON\). To print out the node as-is, with the original formatting an any comments pass in `--unwrapScalar=false`

Given data.yml:

```yaml
a: "Things" # cool stuff
```

Then:

`yq r --unwrapScalar=false data.yml a`

Will yield:

```yaml
"Things" # cool stuff
```

where as without setting the flag to false you would get:

```yaml
Things
```

## Strip comments

Use the `--stripComments` flag to print out the yaml file without any of the original comments.

Given data.yml of:

```yaml
a:
  b: # there is where the good stuff is
    c: hi
```

Then

```yaml
yq r data.yml a --stripComments
```

Will yield:

```yaml
b:
  c: hi
```

