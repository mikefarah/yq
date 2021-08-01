# String Operators

## RegEx
This uses golangs native regex functions under the hood - See https://github.com/google/re2/wiki/Syntax for the supported syntax.


# String blocks, bash and newlines
Bash is notorious for chomping on precious trailing newline characters, making it tricky to set strings with newlines properly. In particular, the `$( exp )` _will trim trailing newlines_.

For instance to get this yaml:

```
a: |
  cat
```

Using `$( exp )` wont work, as it will trim the trailing new line.

```
m=$(echo "cat\n") yq e -n '.a = strenv(m)'
a: cat
```

However, using printf works:
```
printf -v m "cat\n" ; m="$m" yq e -n '.a = strenv(m)'
a: |
  cat
```

As well as having multiline expressions:
```
m="cat
"  yq e -n '.a = strenv(m)'
a: |
  cat
```

Similarly, if you're trying to set the content from a file, and want a trailing new line:

```
IFS= read -rd '' output < <(cat my_file)
output=$output ./yq e '.data.values = strenv(output)' first.yml
```