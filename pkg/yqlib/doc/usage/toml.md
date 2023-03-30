# TOML

Decode from TOML. Note that `yq` does not yet support outputting in TOML format (and therefore it cannot roundtrip)


## Parse: Simple
Given a sample.toml file of:
```toml
A = "hello"
B = 12

```
then
```bash
yq -oy '.' sample.toml
```
will output
```yaml
A: hello
B: 12
```

## Parse: Deep paths
Given a sample.toml file of:
```toml
person.name = "hello"
person.address = "12 cat st"

```
then
```bash
yq -oy '.' sample.toml
```
will output
```yaml
person:
  name: hello
  address: 12 cat st
```

## Parse: inline table
Given a sample.toml file of:
```toml
name = { first = "Tom", last = "Preston-Werner" }
```
then
```bash
yq -oy '.' sample.toml
```
will output
```yaml
name:
  first: Tom
  last: Preston-Werner
```

## Parse: Array Table
Given a sample.toml file of:
```toml

[owner.contact]
name = "Tom Preston-Werner"
age = 36

[[owner.addresses]]
street = "first street"
suburb = "ok"

[[owner.addresses]]
street = "second street"
suburb = "nice"

```
then
```bash
yq -oy '.' sample.toml
```
will output
```yaml
owner:
  contact:
    name: Tom Preston-Werner
    age: 36
  addresses:
    - street: first street
      suburb: ok
    - street: second street
      suburb: nice
```

