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

## Encode: Scalar
Given a sample.toml file of:
```toml
person.name = "hello"
person.address = "12 cat st"

```
then
```bash
yq '.person.name' sample.toml
```
will output
```yaml
hello
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

## Parse: Array of Array Table
Given a sample.toml file of:
```toml

[[fruits]]
name = "apple"
[[fruits.varieties]]  # nested array of tables
name = "red delicious"
```
then
```bash
yq -oy '.' sample.toml
```
will output
```yaml
fruits:
  - name: apple
    varieties:
      - name: red delicious
```

## Parse: Empty Table
Given a sample.toml file of:
```toml

[dependencies]

```
then
```bash
yq -oy '.' sample.toml
```
will output
```yaml
dependencies: {}
```

## Roundtrip: inline table attribute
Given a sample.toml file of:
```toml
name = { first = "Tom", last = "Preston-Werner" }

```
then
```bash
yq '.' sample.toml
```
will output
```yaml
name = { first = "Tom", last = "Preston-Werner" }
```

## Roundtrip: table section
Given a sample.toml file of:
```toml
[owner.contact]
name = "Tom"
age = 36

```
then
```bash
yq '.' sample.toml
```
will output
```yaml
[owner.contact]
name = "Tom"
age = 36
```

## Roundtrip: array of tables
Given a sample.toml file of:
```toml
[[fruits]]
name = "apple"
[[fruits.varieties]]
name = "red delicious"

```
then
```bash
yq '.' sample.toml
```
will output
```yaml
[[fruits]]
name = "apple"
[[fruits.varieties]]
name = "red delicious"
```

## Roundtrip: arrays and scalars
Given a sample.toml file of:
```toml
A = ["hello", ["world", "again"]]
B = 12

```
then
```bash
yq '.' sample.toml
```
will output
```yaml
A = ["hello", ["world", "again"]]
B = 12
```

## Roundtrip: simple
Given a sample.toml file of:
```toml
A = "hello"
B = 12

```
then
```bash
yq '.' sample.toml
```
will output
```yaml
A = "hello"
B = 12
```

## Roundtrip: deep paths
Given a sample.toml file of:
```toml
[person]
name = "hello"
address = "12 cat st"

```
then
```bash
yq '.' sample.toml
```
will output
```yaml
[person]
name = "hello"
address = "12 cat st"
```

## Roundtrip: empty array
Given a sample.toml file of:
```toml
A = []

```
then
```bash
yq '.' sample.toml
```
will output
```yaml
A = []
```

## Roundtrip: sample table
Given a sample.toml file of:
```toml
var = "x"

[owner.contact]
name = "Tom Preston-Werner"
age = 36

```
then
```bash
yq '.' sample.toml
```
will output
```yaml
var = "x"

[owner.contact]
name = "Tom Preston-Werner"
age = 36
```

## Roundtrip: empty table
Given a sample.toml file of:
```toml
[dependencies]

```
then
```bash
yq '.' sample.toml
```
will output
```yaml
[dependencies]
```

## Roundtrip: comments
Given a sample.toml file of:
```toml
# This is a comment
A = "hello"  # inline comment
B = 12

# Table comment
[person]
name = "Tom"  # name comment

```
then
```bash
yq '.' sample.toml
```
will output
```yaml
# This is a comment
A = "hello"  # inline comment
B = 12

# Table comment
[person]
name = "Tom"  # name comment
```

## Roundtrip: sample from web
Given a sample.toml file of:
```toml
# This is a TOML document
title = "TOML Example"

[owner]
name = "Tom Preston-Werner"
dob = 1979-05-27T07:32:00-08:00

[database]
enabled = true
ports = [8000, 8001, 8002]
data = [["delta", "phi"], [3.14]]
temp_targets = { cpu = 79.5, case = 72.0 }

# [servers] yq can't do this one yet
[servers.alpha]
ip = "10.0.0.1"
role = "frontend"

[servers.beta]
ip = "10.0.0.2"
role = "backend"

```
then
```bash
yq '.' sample.toml
```
will output
```yaml
# This is a TOML document
title = "TOML Example"

[owner]
name = "Tom Preston-Werner"
dob = 1979-05-27T07:32:00-08:00

[database]
enabled = true
ports = [8000, 8001, 8002]
data = [["delta", "phi"], [3.14]]
temp_targets = { cpu = 79.5, case = 72.0 }

# [servers] yq can't do this one yet
[servers.alpha]
ip = "10.0.0.1"
role = "frontend"

[servers.beta]
ip = "10.0.0.2"
role = "backend"
```

