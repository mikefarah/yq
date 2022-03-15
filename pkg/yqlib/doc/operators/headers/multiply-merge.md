# Multiply (Merge)

Like the multiple operator in jq, depending on the operands, this multiply operator will do different things. Currently numbers, arrays and objects are supported.

## Objects and arrays - merging
Objects are merged deeply matching on matching keys. By default, array values override and are not deeply merged.

Note that when merging objects, this operator returns the merged object (not the parent). This will be clearer in the examples below.

### Merge Flags
You can control how objects are merged by using one or more of the following flags. Multiple flags can be used together, e.g. `.a *+? .b`.  See examples below

- `+` append arrays
- `d` deeply merge arrays
- `?` only merge _existing_ fields
- `n` only merge _new_ fields


### Merge two files together
This uses the load operator to merge file2 into file1.
```bash
yq '. *= load("file2.yml")' file1.yml
```

### Merging all files
Note the use of `eval-all` to ensure all documents are loaded into memory.

```bash
yq eval-all '. as $item ireduce ({}; . * $item )' *.yml
```
