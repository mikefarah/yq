---
description: Deeply compare two yaml documents
---

# Compare

```bash
yq compare <yaml_file_1> <yaml_file_2> <path_expression>
```

Compares the matching yaml nodes at path expression in the two yaml documents. See [path expression](../usage/path-expressions.md) for more details. Difference calculated line by line, and is printed out line by line where the first character of each line is either:

*  `` a space, indicating no change at this line
* `-` a minus ,indicating the line is not present in the second document \(it's removed\)
* `+` a plus, indicating that the line is not present in the first document \(it's added\)

If there are differences then `yq` will print out the differences and exit with code 1. If there are no differences, then nothing will be printed and the exit code will be 0.

## Example data

Given data1.yaml

```yaml
"apples": are nice
somethingElse: cool # this is nice
favouriteNumbers: [1,2,3]
noDifference: it's the same
```

and data2.yaml

```yaml
apples: are nice
somethingElse: cool # yeah i like it
favouriteNumbers:
- 1
- 3
- 4
noDifference: it's the same
```

## Basic

Basic will compare the yaml documents 'as-is'

```yaml
yq compare data1.yaml data2.yaml
```

yields

```text
-"apples": are nice
-somethingElse: cool # this is nice
-favouriteNumbers: [1, 2, 3]
+apples: are nice
+somethingElse: cool # yeah i like it
+favouriteNumbers:
+- 1
+- 3
+- 4
 noDifference: it's the same
```

## Formatted

Most of the time, it will make sense to [format](../usage/output-format.md#pretty-print) the documents before comparing:

```text
yq compare --prettyPrint data1.yaml data2.yml
```

yields

```text
 apples: are nice
-somethingElse: cool # this is nice
+somethingElse: cool # yeah i like it
 favouriteNumbers:
 - 1
-- 2
 - 3
+- 4
 noDifference: it's the same
```

## Using path expressions

Use [path expressions](../usage/path-expressions.md) to compare subsets of yaml documents

```text
yq compare -P data1.yaml data2.yml favouriteNumbers
```

yields

```text
 - 1
-- 2
 - 3
+- 4
```

