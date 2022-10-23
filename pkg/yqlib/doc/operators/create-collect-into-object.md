# Create, Collect into Object

This is used to construct objects (or maps). This can be used against existing yaml, or to create fresh yaml documents.

## Collect empty object
Running
```bash
yq --null-input '{}'
```
will output
```yaml
{}
```

## Wrap (prefix) existing object
Given a sample.yml file of:
```yaml
name: Mike
```
then
```bash
yq '{"wrap": .}' sample.yml
```
will output
```yaml
wrap:
  name: Mike
```

## Using splat to create multiple objects
Given a sample.yml file of:
```yaml
name: Mike
pets:
  - cat
  - dog
```
then
```bash
yq '{.name: .pets.[]}' sample.yml
```
will output
```yaml
Mike: cat
Mike: dog
```

## Working with multiple documents
Given a sample.yml file of:
```yaml
name: Mike
pets:
  - cat
  - dog
---
name: Rosey
pets:
  - monkey
  - sheep
```
then
```bash
yq '{.name: .pets.[]}' sample.yml
```
will output
```yaml
Mike: cat
Mike: dog
Rosey: monkey
Rosey: sheep
```

## Creating yaml from scratch
Running
```bash
yq --null-input '{"wrap": "frog"}'
```
will output
```yaml
wrap: frog
```

## Creating yaml from scratch with multiple objects
Running
```bash
yq --null-input '(.a.b = "foo") | (.d.e = "bar")'
```
will output
```yaml
a:
  b: foo
d:
  e: bar
```

