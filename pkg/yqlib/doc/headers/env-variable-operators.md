# Env Variable Operators

This operator is used to handle environment variables usage in path expressions. While environment variables can, of course, be passed in via your CLI with string interpolation, this often comes with complex quote escaping and can be tricky to write and read. Note that there are two forms, `env` which will parse the environment variable as a yaml (be it a map, array, string, number of boolean) and `strenv` which will always parse the argument as a string.
