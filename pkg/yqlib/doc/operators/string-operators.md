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

## To up (upper) case
Works with unicode characters

Given a sample.yml file of:
```yaml
água
```
then
```bash
yq 'upcase' sample.yml
```
will output
```yaml
ÁGUA
```

## To down (lower) case
Works with unicode characters

Given a sample.yml file of:
```yaml
ÁgUA
```
then
```bash
yq 'downcase' sample.yml
```
will output
```yaml
água
```

## Join strings
Given a sample.yml file of:
```yaml
- cat
- meow
- 1
- null
- true
```
then
```bash
yq 'join("; ")' sample.yml
```
will output
```yaml
cat; meow; 1; ; true
```

## Trim strings
Given a sample.yml file of:
```yaml
- ' cat'
- 'dog '
- ' cow cow '
- horse
```
then
```bash
yq '.[] | trim' sample.yml
```
will output
```yaml
cat
dog
cow cow
horse
```

## Match string
Given a sample.yml file of:
```yaml
foo bar foo
```
then
```bash
yq 'match("foo")' sample.yml
```
will output
```yaml
string: foo
offset: 0
length: 3
captures: []
```

## Match string, case insensitive
Given a sample.yml file of:
```yaml
foo bar FOO
```
then
```bash
yq '[match("(?i)foo"; "g")]' sample.yml
```
will output
```yaml
- string: foo
  offset: 0
  length: 3
  captures: []
- string: FOO
  offset: 8
  length: 3
  captures: []
```

## Match with global capture group
Given a sample.yml file of:
```yaml
abc abc
```
then
```bash
yq '[match("(ab)(c)"; "g")]' sample.yml
```
will output
```yaml
- string: abc
  offset: 0
  length: 3
  captures:
    - string: ab
      offset: 0
      length: 2
    - string: c
      offset: 2
      length: 1
- string: abc
  offset: 4
  length: 3
  captures:
    - string: ab
      offset: 4
      length: 2
    - string: c
      offset: 6
      length: 1
```

## Match with named capture groups
Given a sample.yml file of:
```yaml
foo bar foo foo  foo
```
then
```bash
yq '[match("foo (?P<bar123>bar)? foo"; "g")]' sample.yml
```
will output
```yaml
- string: foo bar foo
  offset: 0
  length: 11
  captures:
    - string: bar
      offset: 4
      length: 3
      name: bar123
- string: foo  foo
  offset: 12
  length: 8
  captures:
    - string: null
      offset: -1
      length: 0
      name: bar123
```

## Capture named groups into a map
Given a sample.yml file of:
```yaml
xyzzy-14
```
then
```bash
yq 'capture("(?P<a>[a-z]+)-(?P<n>[0-9]+)")' sample.yml
```
will output
```yaml
a: xyzzy
n: "14"
```

## Match without global flag
Given a sample.yml file of:
```yaml
cat cat
```
then
```bash
yq 'match("cat")' sample.yml
```
will output
```yaml
string: cat
offset: 0
length: 3
captures: []
```

## Match with global flag
Given a sample.yml file of:
```yaml
cat cat
```
then
```bash
yq '[match("cat"; "g")]' sample.yml
```
will output
```yaml
- string: cat
  offset: 0
  length: 3
  captures: []
- string: cat
  offset: 4
  length: 3
  captures: []
```

## Test using regex
Like jq's equivalent, this works like match but only returns true/false instead of full match details

Given a sample.yml file of:
```yaml
- cat
- dog
```
then
```bash
yq '.[] | test("at")' sample.yml
```
will output
```yaml
true
false
```

## Substitute / Replace string
This uses Golang's regex, described [here](https://github.com/google/re2/wiki/Syntax)
Note the use of `|=` to run in context of the current string value.

Given a sample.yml file of:
```yaml
a: dogs are great
```
then
```bash
yq '.a |= sub("dogs", "cats")' sample.yml
```
will output
```yaml
a: cats are great
```

## Substitute / Replace string with regex
This uses Golang's regex, described [here](https://github.com/google/re2/wiki/Syntax)
Note the use of `|=` to run in context of the current string value.

Given a sample.yml file of:
```yaml
a: cat
b: heat
```
then
```bash
yq '.[] |= sub("(a)", "${1}r")' sample.yml
```
will output
```yaml
a: cart
b: heart
```

## Custom types: that are really strings
When custom tags are encountered, yq will try to decode the underlying type.

Given a sample.yml file of:
```yaml
a: !horse cat
b: !goat heat
```
then
```bash
yq '.[] |= sub("(a)", "${1}r")' sample.yml
```
will output
```yaml
a: !horse cart
b: !goat heart
```

## Split strings
Given a sample.yml file of:
```yaml
cat; meow; 1; ; true
```
then
```bash
yq 'split("; ")' sample.yml
```
will output
```yaml
- cat
- meow
- "1"
- ""
- "true"
```

## Split strings one match
Given a sample.yml file of:
```yaml
word
```
then
```bash
yq 'split("; ")' sample.yml
```
will output
```yaml
- word
```

