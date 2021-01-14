Like the multiple operator in jq, depending on the operands, this multiply operator will do different things. Currently only objects are supported, which have the effect of merging the RHS into the LHS.

To concatenate arrays when merging objects, use the *+ form (see examples below). This will recursively merge objects, appending arrays when it encounters them.

To merge only existing fields, use the *? form. Note that this can be used with the concatenate arrays too *+?.
Note that when merging objects, this operator returns the merged object (not the parent). This will be clearer in the examples below.

Multiplication of strings and numbers are not yet supported.

## Merging files
Note the use of `eval-all` to ensure all documents are loaded into memory.

```bash
yq eval-all 'select(fileIndex == 0) * select(fileIndex == 1)' file1.yaml file2.yaml
```
