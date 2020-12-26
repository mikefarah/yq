Use these comment operators to set or retrieve comments.
## Set line comment
Given a sample.yml file of:
```yaml
a: cat
'': null
```
then
```bash
yq eval '.a lineComment="single"' sample.yml
```
will output
```yaml
a: cat # single
'': null
```

## Set head comment
Given a sample.yml file of:
```yaml
a: cat
'': null
```
then
```bash
yq eval '. headComment="single"' sample.yml
```
will output
```yaml
# single

a: cat
'': null
```

## Set foot comment, using an expression
Given a sample.yml file of:
```yaml
a: cat
'': null
```
then
```bash
yq eval '. footComment=.a' sample.yml
```
will output
```yaml
a: cat
'': null

# cat
```

## Remove comment
Given a sample.yml file of:
```yaml
a: cat # comment
b: dog # leave this
'': null
```
then
```bash
yq eval '.a lineComment=""' sample.yml
```
will output
```yaml
a: cat
b: dog # leave this
'': null
```

## Remove all comments
Given a sample.yml file of:
```yaml
a: cat # comment
'': null
```
then
```bash
yq eval '.. comments=""' sample.yml
```
will output
```yaml
a: cat # comment
'': null
```

## Get line comment
Given a sample.yml file of:
```yaml
a: cat # meow
'': null
```
then
```bash
yq eval '.a | lineComment' sample.yml
```
will output
```yaml
meow
```

## Get head comment
Given a sample.yml file of:
```yaml
a: cat # meow
'': null
```
then
```bash
yq eval '. | headComment' sample.yml
```
will output
```yaml

```

## Get foot comment
Given a sample.yml file of:
```yaml
a: cat # meow
'': null
```
then
```bash
yq eval '. | footComment' sample.yml
```
will output
```yaml

```

