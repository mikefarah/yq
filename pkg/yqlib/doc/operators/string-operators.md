# String Operators

## RegEx
This uses golangs native regex functions under the hood - See https://github.com/google/re2/wiki/Syntax for the supported syntax.


## String blocks, bash and newlines
Bash is notorious for chomping on precious trailing newline characters, making it tricky to set strings with newlines properly. In particular, the `$( exp )` _will trim trailing newlines_.

For instance to get this yaml:

```
a: |
  cat
```

Using `$( exp )` wont work, as it will trim the trailing new line.

```
m=$(echo "cat\n") yq e -n '.a = strenv(m)'
a: cat
```

However, using printf works:
```
printf -v m "cat\n" ; m="$m" yq e -n '.a = strenv(m)'
a: |
  cat
```

As well as having multiline expressions:
```
m="cat
"  yq e -n '.a = strenv(m)'
a: |
  cat
```

Similarly, if you're trying to set the content from a file, and want a trailing new line:

```
IFS= read -rd '' output < <(cat my_file)
output=$output ./yq e '.data.values = strenv(output)' first.yml
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
yq eval 'join("; ")' sample.yml
```
will output
```yaml
cat; meow; 1; ; true
```

## Match string
Given a sample.yml file of:
```yaml
foo bar foo
```
then
```bash
yq eval 'match("foo")' sample.yml
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
yq eval '[match("(?i)foo"; "g")]' sample.yml
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

## Match with capture groups
Given a sample.yml file of:
```yaml
abc abc
```
then
```bash
yq eval '[match("(abc)+"; "g")]' sample.yml
```
will output
```yaml
- string: abc
  offset: 0
  length: 3
  captures:
    - string: abc
      offset: 0
      length: 3
- string: abc
  offset: 4
  length: 3
  captures:
    - string: abc
      offset: 4
      length: 3
```

## Match with named capture groups
Given a sample.yml file of:
```yaml
foo bar foo foo  foo
```
then
```bash
yq eval '[match("foo (?P<bar123>bar)? foo"; "g")]' sample.yml
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
yq eval 'capture("(?P<a>[a-z]+)-(?P<n>[0-9]+)")' sample.yml
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
yq eval 'match("cat")' sample.yml
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
yq eval '[match("cat"; "g")]' sample.yml
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
Like jq'q equivalent, this works like match but only returns true/false instead of full match details

Given a sample.yml file of:
```yaml
- cat
- dog
```
then
```bash
yq eval '.[] | test("at")' sample.yml
```
will output
```yaml
true
false
```

## Substitute / Replace string
This uses golang regex, described [here](https://github.com/google/re2/wiki/Syntax)
Note the use of `|=` to run in context of the current string value.

Given a sample.yml file of:
```yaml
a: dogs are great
```
then
```bash
yq eval '.a |= sub("dogs", "cats")' sample.yml
```
will output
```yaml
a: cats are great
```

## Substitute / Replace string with regex
This uses golang regex, described [here](https://github.com/google/re2/wiki/Syntax)
Note the use of `|=` to run in context of the current string value.

Given a sample.yml file of:
```yaml
a: cat
b: heat
```
then
```bash
yq eval '.[] |= sub("(a)", "${1}r")' sample.yml
```
will output
```yaml
a: cart
b: heart
```

## Split strings
Given a sample.yml file of:
```yaml
cat; meow; 1; ; true
```
then
```bash
yq eval 'split("; ")' sample.yml
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
yq eval 'split("; ")' sample.yml
```
will output
```yaml
- word
```

