# String Operators

## RegEx
This uses golangs native regex functions under the hood - See https://github.com/google/re2/wiki/Syntax for the supported syntax.

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
cat
```
then
```bash
yq eval 'match("at")' sample.yml
```
will output
```yaml
string: at
offset: 1
length: 2
captures: []
```

## Match string, case insensitive
Given a sample.yml file of:
```yaml
cAt
```
then
```bash
yq eval 'match("(?i)at")' sample.yml
```
will output
```yaml
string: At
offset: 1
length: 2
captures: []
```

## Match with capture groups
Given a sample.yml file of:
```yaml
a cat
```
then
```bash
yq eval 'match("c(.t)")' sample.yml
```
will output
```yaml
string: cat
offset: 2
length: 3
captures:
  - string: at
    offset: 3
    length: 2
```

## Match with named capture groups
Given a sample.yml file of:
```yaml
a cat
```
then
```bash
yq eval 'match("c(?P<cool>.t)")' sample.yml
```
will output
```yaml
string: cat
offset: 2
length: 3
captures:
  - string: at
    offset: 3
    length: 2
    name: cool
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
yq eval 'match("cat"; "g")' sample.yml
```
will output
```yaml
string: cat
offset: 0
length: 3
captures: []
string: cat
offset: 4
length: 3
captures: []
```

## Test using regex
Like jq'q equivalant, this works like match but only returns true/false instead of full match details

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

