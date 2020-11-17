Use these comment operators to set or retrieve comments.
## Examples
### Set line comment
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq eval '.a lineComment="single"' sample.yml
```
will output
```yaml
a: cat # single
```

### Set head comment
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq eval '. headComment="single"' sample.yml
```
will output
```yaml
# single

a: cat
```

### Set foot comment, using an expression
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq eval '. footComment=.a' sample.yml
```
will output
```yaml
a: cat

# cat
```

### Remove comment
Given a sample.yml file of:
```yaml
a: cat # comment
b: dog # leave this
```
then
```bash
yq eval '.a lineComment=""' sample.yml
```
will output
```yaml
a: cat
b: dog # leave this
```

### Remove all comments
Given a sample.yml file of:
```yaml
a: cat # comment
```
then
```bash
yq eval '.. comments=""' sample.yml
```
will output
```yaml
a: cat
```

### Get line comment
Given a sample.yml file of:
```yaml
a: cat # meow
```
then
```bash
yq eval '.a | lineComment' sample.yml
```
will output
```yaml
meow
```

### Get head comment
Given a sample.yml file of:
```yaml
a: cat # meow
```
then
```bash
yq eval '. | headComment' sample.yml
```
will output
```yaml
welcome!
```

### Get foot comment
Given a sample.yml file of:
```yaml
a: cat # meow
```
then
```bash
yq eval '. | footComment' sample.yml
```
will output
```yaml
have a great day
```

