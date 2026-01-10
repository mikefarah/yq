# JSON5

JSON5 support in `yq` lets you parse JSON5 files (comments, trailing commas, single quotes, unquoted keys, hex numbers, `Infinity`, `NaN`) and convert them to other formats like YAML, or output JSON5.

Note: when converting JSON5 to YAML (or other formats), comments may move slightly because formats like YAML don't always have a distinct representation for certain JSON5 comment placements (e.g. `/* foo */ { ... }` vs `{ /* foo */ ... }`). When converting JSON5 back to JSON5, `yq` keeps comments as close as possible to their original location.

## Parse json5: comments, trailing commas, single quotes
Given a sample.json5 file of:
```json5
{
  // comment
  unquoted: 'single quoted',
  trailing: [1, 2,],
}

```
then
```bash
yq -P -p=json5 '.' sample.json5
```
will output
```yaml
# comment
unquoted: single quoted
trailing:
  - 1
  - 2
```

