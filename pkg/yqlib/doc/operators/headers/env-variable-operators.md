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

