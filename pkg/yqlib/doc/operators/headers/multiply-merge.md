# Multiply (Merge)

Like the multiple operator in jq, depending on the operands, this multiply operator will do different things. Currently numbers, arrays and objects are supported.

## Objects and arrays - merging
Objects are merged _deeply_ matching on matching keys. By default, array values override and are not deeply merged.

You can use the add operator `+`, to shallow merge objects, see more info [here](https://mikefarah.gitbook.io/yq/operators/add).

Note that when merging objects, this operator returns the merged object (not the parent). This will be clearer in the examples below.

### Merge Flags
You can control how objects are merged by using one or more of the following flags. Multiple flags can be used together, e.g. `.a *+? .b`.  See examples below

- `+` append arrays
- `d` deeply merge arrays
- `?` only merge _existing_ fields
- `n` only merge _new_ fields
- `c` clobber custom tags

To perform a shallow merge only, use the add operator `+`, see more info [here](https://mikefarah.gitbook.io/yq/operators/add).

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

# Merging complex arrays together by a key field
By default - `yq` merge is naive. It merges maps when they match the key name, and arrays are merged either by appending them together, or merging the entries by their position in the array.

For more complex array merging (e.g. merging items that match on a certain key) please see the example [here](https://mikefarah.gitbook.io/yq/operators/multiply-merge#merge-arrays-of-objects-together-matching-on-a-key)

