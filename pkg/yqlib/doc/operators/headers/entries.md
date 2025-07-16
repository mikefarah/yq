# Entries

Similar to the same named functions in `jq` these functions convert to/from an object and an array of key-value pairs. This is most useful for performing operations on keys of maps.

Use `with_entries(op)` as a syntactic sugar for doing `to_entries | op | from_entries`.
