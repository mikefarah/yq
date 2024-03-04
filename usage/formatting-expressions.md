# Formatting Expressions

`From version v4.41+`

You can put expressions into `.yq` files, use whitespace and comments to break up complex expressions and explain what's going on.

## Using expression files and comments
Note that you can execute the file directly - but make sure you make the expression file executable.

Given a sample.yaml file of:
```yaml
a:
  b: old
```
And an 'update.yq' expression file of:
```bash
#! yq

# This is a yq expression that updates the map
# for several great reasons outlined here.

.a.b = "new" # line comment here
| .a.c = "frog"

# Now good things will happen.
```
then
```bash
./update.yq sample.yaml
```
will output
```yaml
a:
  b: new
  c: frog
```

## Flags in expression files
You can specify flags on the shebang line, this only works when executing the file directly.

Given a sample.yaml file of:
```yaml
a:
  b: old
```
And an 'update.yq' expression file of:
```bash
#! yq -oj

# This is a yq expression that updates the map
# for several great reasons outlined here.

.a.b = "new" # line comment here
| .a.c = "frog"

# Now good things will happen.
```
then
```bash
./update.yq sample.yaml
```
will output
```yaml
{
  "a": {
    "b": "new",
    "c": "frog"
  }
}
```

## Commenting out yq expressions
Note that `c` is no longer set to 'frog'. In this example we're calling yq directly and passing the expression file into `--from-file`, this is no different from executing the expression file directly.

Given a sample.yaml file of:
```yaml
a:
  b: old
```
And an 'update.yq' expression file of:
```bash
#! yq
# This is a yq expression that updates the map
# for several great reasons outlined here.

.a.b = "new" # line comment here
# | .a.c = "frog"

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
```

