# Array to Map

Use this operator to convert an array to..a map. Skips over null values.

Behind the scenes, this is implemented using reduce:

```
(.[] | select(. != null) ) as $i ireduce({}; .[$i | key] = $i)
```
