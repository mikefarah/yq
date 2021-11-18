# Load

The `load`/`strload` operator allows you to load in content from another file referenced in your yaml document.

Note that you can use string operators like `+` and `sub` to modify the value in the yaml file to a path that exists in your system.

Use `strload` to load text based content as a string block, and `load` to interpret the file as yaml.

Lets say there is a file `../../examples/thing.yml`:

```yaml
a: apple is included
b: cool
```
