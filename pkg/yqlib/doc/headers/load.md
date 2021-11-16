# Load

The `load`/`strload` operator allows you to load in content from another file referenced in your yaml document.

Note that you can use string operators like `+` and `sub` to modify the value in the yaml file to a path that exists in your system.


Lets say there is a file `../../examples/thing.yml`:

```yaml
a: apple is included
b: cool
```
