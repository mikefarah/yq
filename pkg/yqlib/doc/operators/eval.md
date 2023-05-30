# Eval

Use `eval` to dynamically process an expression - for instance from an environment variable.

`eval` takes a single argument, and evaluates that as a `yq` expression. Any valid expression can be used, be it a path `.a.b.c | select(. == "cat")`, or an update `.a.b.c = "gogo"`.

Tip: This can be a useful way to parameterise complex scripts.

## Dynamically evaluate a path
Given a sample.yml file of:
```yaml
pathExp: .a.b[] | select(.name == "cat")
a:
  b:
    - name: dog
    - name: cat
```
then
```bash
yq 'eval(.pathExp)' sample.yml
```
will output
```yaml
name: cat
```

## Dynamically update a path from an environment variable
The env variable can be any valid yq expression.

Given a sample.yml file of:
```yaml
a:
  b:
    - name: dog
    - name: cat
```
then
```bash
pathEnv=".a.b[0].name"  valueEnv="moo" yq 'eval(strenv(pathEnv)) = strenv(valueEnv)' sample.yml
```
will output
```yaml
a:
  b:
    - name: moo
    - name: cat
```

