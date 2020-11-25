Add behaves differently according to the type of the LHS:
- arrays: concatenate
- number scalars: arithmetic addition (soon)
- string scalars: concatenate (soon)

## Concatenate arrays
Given a sample.yml file of:
```yaml
a:
  - 1
  - 2
b:
  - 3
  - 4
```
then
```bash
yq eval '.a + .b' sample.yml
```
will output
```yaml
- 1
- 2
- 3
- 4
```

## Concatenate null to array
Given a sample.yml file of:
```yaml
a:
  - 1
  - 2
```
then
```bash
yq eval '.a + null' sample.yml
```
will output
```yaml
- 1
- 2
```

## Add object to array
Given a sample.yml file of:
```yaml
a:
  - 1
  - 2
c:
  cat: meow
```
then
```bash
yq eval '.a + .c' sample.yml
```
will output
```yaml
- 1
- 2
- cat: meow
```

## Add string to array
Given a sample.yml file of:
```yaml
a:
  - 1
  - 2
```
then
```bash
yq eval '.a + "hello"' sample.yml
```
will output
```yaml
- 1
- 2
- hello
```

## Update array (append)
Given a sample.yml file of:
```yaml
a:
  - 1
  - 2
b:
  - 3
  - 4
```
then
```bash
yq eval '.a = .a + .b' sample.yml
```
will output
```yaml
a:
  - 1
  - 2
  - 3
  - 4
b:
  - 3
  - 4
```

