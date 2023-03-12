# Sort Keys

The Sort Keys operator sorts maps by their keys (based on their string value). This operator does not do anything to arrays or scalars (so you can easily recursively apply it to all maps).

Sort is particularly useful for diffing two different yaml documents:

```bash
yq -i -P 'sort_keys(..)' file1.yml
yq -i -P 'sort_keys(..)' file2.yml
diff file1.yml file2.yml
```

Note that `yq` does not yet consider anchors when sorting by keys - this may result in invalid yaml documents if you are using merge anchors.

For more advanced sorting, using `to_entries` to convert the map to an array, then sort/process the array as you like (e.g. using `sort_by`) and convert back to a map using `from_entries`.
See [here](https://mikefarah.gitbook.io/yq/operators/entries#custom-sort-map-keys) for an example. 
