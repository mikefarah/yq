# Formatting Expressions

`From version v4.41+`

You can put expressions into `.yq` files, use whitespace and comments to break up complex expressions and explain what's going on.

## Using expression files and comments
Given a sample.yaml file of:
```yaml
a:
  b: old
```
And an 'update.yq' expression file of:
```bash
# This is a yq expression that updates the map
# for several great reasons outlined here.

.a.b = "new" # line comment here
| .a.c = "frog"

# Now good things will happen.
```
then
```bash
yq --from-file update.yq sample.yml
```
will output
```yaml
a:
  b: new
  c: frog
```

