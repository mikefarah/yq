# Envsubst

This operator is used to replace environment variables in strings using [envsubst](https://github.com/a8m/envsubst).

To replace environment variables across all values in a document, this can be used with the recursive descent operator
as follows:

```bash
yq eval '(.. | select(tag == "!!str")) |= envsubst' file.yaml
```
