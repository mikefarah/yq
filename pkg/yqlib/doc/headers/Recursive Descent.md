This operator recursively matches (or globs) all children nodes given of a particular element, including that node itself. This is most often used to apply a filter recursively against all matches. It can be used in either the 

## match values form `..`
This will, like the `jq` equivalent, recursively match all _value_ nodes. Use it to find/manipulate particular values.

For instance to set the `style` of all _value_ nodes in a yaml doc, excluding map keys:

```bash
yq eval '.. style= "flow"' file.yaml
```

## match values and map keys form `...`
The also includes map keys in the results set. This is particularly useful in YAML as unlike JSON, map keys can have their own styling, tags and use anchors and aliases.

For instance to set the `style` of all nodes in a yaml doc, including the map keys:

```bash
yq eval '... style= "flow"' file.yaml
```