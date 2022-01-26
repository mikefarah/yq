# Env Variable Operators

These operators are used to handle environment variables usage in expressions and documents. While environment variables can, of course, be passed in via your CLI with string interpolation, this often comes with complex quote escaping and can be tricky to write and read. 

There are three operators:
-  `env` which takes a single environment variable name and parse the variable as a yaml node (be it a map, array, string, number of boolean) 
- `strenv` which also takes a single environment variable name, and always parses the variable as a string.
- `envsubst` which you pipe strings into and it interpolates environment variables in strings using [envsubst](https://github.com/a8m/envsubst). 


## Tip
To replace environment variables across all values in a document, `envsubst` can be used with the recursive descent operator
as follows:

```bash
yq eval '(.. | select(tag == "!!str")) |= envsubst' file.yaml
```

