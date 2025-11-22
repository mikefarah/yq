# Env Variable Operators

These operators are used to handle environment variables usage in expressions and documents. While environment variables can, of course, be passed in via your CLI with string interpolation, this often comes with complex quote escaping and can be tricky to write and read. 

There are three operators:

-  `env` which takes a single environment variable name and parse the variable as a yaml node (be it a map, array, string, number of boolean) 
- `strenv` which also takes a single environment variable name, and always parses the variable as a string.
- `envsubst` which you pipe strings into and it interpolates environment variables in strings using [envsubst](https://github.com/a8m/envsubst). 


## EnvSubst Options
You can optionally pass envsubst any of the following options:

  - nu: NoUnset, this will fail if there are any referenced variables that are not set
  - ne: NoEmpty, this will fail if there are any referenced variables that are empty
  - ff: FailFast, this will abort on the first failure (rather than collect all the errors)

E.g:
`envsubst(ne, ff)` will fail on the first empty variable.

See [Imposing Restrictions](https://github.com/a8m/envsubst#imposing-restrictions) in the `envsubst` documentation for more information, and below for examples.

## Tip
To replace environment variables across all values in a document, `envsubst` can be used with the recursive descent operator
as follows:

```bash
yq '(.. | select(tag == "!!str")) |= envsubst' file.yaml
```

## Disabling env operators
If required, you can use the `--security-disable-env-ops` to disable env operations.


## Read string environment variable
Running
```bash
myenv="cat meow" yq --null-input '.a = env(myenv)'
```
will output
```yaml
a: cat meow
```

## Read boolean environment variable
Running
```bash
myenv="true" yq --null-input '.a = env(myenv)'
```
will output
```yaml
a: true
```

## Read numeric environment variable
Running
```bash
myenv="12" yq --null-input '.a = env(myenv)'
```
will output
```yaml
a: 12
```

## Read yaml environment variable
Running
```bash
myenv="{b: fish}" yq --null-input '.a = env(myenv)'
```
will output
```yaml
a: {b: fish}
```

## Read boolean environment variable as a string
Running
```bash
myenv="true" yq --null-input '.a = strenv(myenv)'
```
will output
```yaml
a: "true"
```

## Read numeric environment variable as a string
Running
```bash
myenv="12" yq --null-input '.a = strenv(myenv)'
```
will output
```yaml
a: "12"
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

## Dynamic key lookup with environment variable
Given a sample.yml file of:
```yaml
cat: meow
dog: woof
```
then
```bash
myenv="cat" yq '.[env(myenv)]' sample.yml
```
will output
```yaml
meow
```

## Replace strings with envsubst
Running
```bash
myenv="cat" yq --null-input '"the ${myenv} meows" | envsubst'
```
will output
```yaml
the cat meows
```

## Replace strings with envsubst, missing variables
Running
```bash
yq --null-input '"the ${myenvnonexisting} meows" | envsubst'
```
will output
```yaml
the  meows
```

## Replace strings with envsubst(nu), missing variables
(nu) not unset, will fail if there are unset (missing) variables

Running
```bash
yq --null-input '"the ${myenvnonexisting} meows" | envsubst(nu)'
```
will output
```bash
Error: variable ${myenvnonexisting} not set
```

## Replace strings with envsubst(ne), missing variables
(ne) not empty, only validates set variables

Running
```bash
yq --null-input '"the ${myenvnonexisting} meows" | envsubst(ne)'
```
will output
```yaml
the  meows
```

## Replace strings with envsubst(ne), empty variable
(ne) not empty, will fail if a references variable is empty

Running
```bash
myenv="" yq --null-input '"the ${myenv} meows" | envsubst(ne)'
```
will output
```bash
Error: variable ${myenv} set but empty
```

## Replace strings with envsubst, missing variables with defaults
Running
```bash
yq --null-input '"the ${myenvnonexisting-dog} meows" | envsubst'
```
will output
```yaml
the dog meows
```

## Replace strings with envsubst(nu), missing variables with defaults
Having a default specified skips over the missing variable.

Running
```bash
yq --null-input '"the ${myenvnonexisting-dog} meows" | envsubst(nu)'
```
will output
```yaml
the dog meows
```

## Replace strings with envsubst(ne), missing variables with defaults
Fails, because the variable is explicitly set to blank.

Running
```bash
myEmptyEnv="" yq --null-input '"the ${myEmptyEnv-dog} meows" | envsubst(ne)'
```
will output
```bash
Error: variable ${myEmptyEnv} set but empty
```

## Replace string environment variable in document
Given a sample.yml file of:
```yaml
v: ${myenv}
```
then
```bash
myenv="cat meow" yq '.v |= envsubst' sample.yml
```
will output
```yaml
v: cat meow
```

## (Default) Return all envsubst errors
By default, all errors are returned at once.

Running
```bash
yq --null-input '"the ${notThere} ${alsoNotThere}" | envsubst(nu)'
```
will output
```bash
Error: variable ${notThere} not set
variable ${alsoNotThere} not set
```

## Fail fast, return the first envsubst error (and abort)
Running
```bash
yq --null-input '"the ${notThere} ${alsoNotThere}" | envsubst(nu,ff)'
```
will output
```bash
Error: variable ${notThere} not set
```

## env() operation fails when security is enabled
Use `--security-disable-env-ops` to disable env operations for security.

Running
```bash
yq --null-input 'env("MYENV")'
```
will output
```bash
Error: env operations have been disabled
```

## strenv() operation fails when security is enabled
Use `--security-disable-env-ops` to disable env operations for security.

Running
```bash
yq --null-input 'strenv("MYENV")'
```
will output
```bash
Error: env operations have been disabled
```

## envsubst() operation fails when security is enabled
Use `--security-disable-env-ops` to disable env operations for security.

Running
```bash
yq --null-input '"value: ${MYENV}" | envsubst'
```
will output
```bash
Error: env operations have been disabled
```

