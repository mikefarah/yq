Like the multiple operator in `jq`, depending on the operands, this multiply operator will do different things. Currently only objects are supported, which have the effect of merging the RHS into the LHS.

Upcoming versions of `yq` will add support for other types of multiplication (numbers, strings).

Note that when merging objects, this operator returns the merged object (not the parent). This will be clearer in the examples below.

## Merging files
Note the use of eval-all to ensure all documents are loaded into memory.

```bash
yq eval-all 'select(fileIndex == 0) * select(fileIndex == 1)' file1.yaml file2.yaml
```
