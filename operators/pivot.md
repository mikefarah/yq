# Pivot

Emulates the `PIVOT` function supported by several popular RDBMS systems.

## Pivot a sequence of sequences
Given a sample.yml file of:
```yaml
- - foo
  - bar
  - baz
- - sis
  - boom
  - bah
```
then
```bash
yq 'pivot' sample.yml
```
will output
```yaml
- - foo
  - sis
- - bar
  - boom
- - baz
  - bah
```

## Pivot sequence of heterogeneous sequences
Missing values are "padded" to null.

Given a sample.yml file of:
```yaml
- - foo
  - bar
  - baz
- - sis
  - boom
  - bah
  - blah
```
then
```bash
yq 'pivot' sample.yml
```
will output
```yaml
- - foo
  - sis
- - bar
  - boom
- - baz
  - bah
- -
  - blah
```

## Pivot sequence of maps
Given a sample.yml file of:
```yaml
- foo: a
  bar: b
  baz: c
- foo: x
  bar: y
  baz: z
```
then
```bash
yq 'pivot' sample.yml
```
will output
```yaml
foo:
  - a
  - x
bar:
  - b
  - y
baz:
  - c
  - z
```

## Pivot sequence of heterogeneous maps
Missing values are "padded" to null.

Given a sample.yml file of:
```yaml
- foo: a
  bar: b
  baz: c
- foo: x
  bar: y
  baz: z
  what: ever
```
then
```bash
yq 'pivot' sample.yml
```
will output
```yaml
foo:
  - a
  - x
bar:
  - b
  - y
baz:
  - c
  - z
what:
  -
  - ever
```

