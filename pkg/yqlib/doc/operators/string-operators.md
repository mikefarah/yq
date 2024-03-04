# String Operators

## RegEx
This uses Golang's native regex functions under the hood - See their [docs](https://github.com/google/re2/wiki/Syntax) for the supported syntax.

Case insensitive tip: prefix the regex with `(?i)` - e.g. `test("(?i)cats)"`.

### match(regEx)
This operator returns the substring match details of the given regEx.

### capture(regEx)
Capture returns named RegEx capture groups in a map. Can be more convenient than `match` depending on what you are doing.

## test(regEx)
Returns true if the string matches the RegEx, false otherwise.

## sub(regEx, replacement)
Substitutes matched substrings. The first parameter is the regEx to match substrings within the original string. The second parameter specifies what to replace those matches with. This can refer to capture groups from the first RegEx.

## String blocks, bash and newlines
Bash is notorious for chomping on precious trailing newline characters, making it tricky to set strings with newlines properly. In particular, the `$( exp )` _will trim trailing newlines_.

For instance to get this yaml:

```
a: |
  cat
```

Using `$( exp )` wont work, as it will trim the trailing newline.

```
m=$(echo "cat\n") yq -n '.a = strenv(m)'
a: cat
```

However, using printf works:
```
printf -v m "cat\n" ; m="$m" yq -n '.a = strenv(m)'
a: |
  cat
```

As well as having multiline expressions:
```
m="cat
"  yq -n '.a = strenv(m)'
a: |
  cat
```

Similarly, if you're trying to set the content from a file, and want a trailing newline:

```
IFS= read -rd '' output < <(cat my_file)
output=$output ./yq '.data.values = strenv(output)' first.yml
```

## To string
Note that you may want to force `yq` to leave scalar values wrapped by passing in `--unwrapScalar=false` or `-r=f`

Given a sample.yml file of:
```yaml
- 1
- true
- null
- ~
- cat
- an: object
- - array
  - 2
```
then
```bash
yq '.[] |= to_string' sample.yml
```
will output
```yaml
- "1"
- "true"
- "null"
- "~"
- cat
- "an: object"
- "- array\n- 2"
```

