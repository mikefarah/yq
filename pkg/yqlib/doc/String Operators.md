# String Operators

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

## Substitute / Replace string
This uses golang regex, described [here](https://github.com/google/re2/wiki/Syntax)

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

Given a sample.yml file of:
```yaml
a: cat
b: heat
```
then
```bash
yq eval '.[] |= sub("([a])", "${1}r")' sample.yml
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

