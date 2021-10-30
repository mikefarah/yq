# Union

This operator is used to combine different results together.

## Combine scalars

Running

```bash
yq eval --null-input '1, true, "cat"'
```

will output

```yaml
1
true
cat
```

## Combine selected paths

Given a sample.yml file of:

```yaml
a: fieldA
b: fieldB
c: fieldC
```

then

```bash
yq eval '.a, .c' sample.yml
```

will output

```yaml
fieldA
fieldC
```
