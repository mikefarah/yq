# With

Use the `with` operator to conveniently make multiple updates to a deeply nested path, or to update array elements relatively to each other. The first argument expression sets the root context, and the second expression runs against that root context.

## Update and style
Given a sample.yml file of:
```yaml
a:
  deeply:
    nested: value
```
then
```bash
yq eval 'with(.a.deeply.nested; . = "newValue" | . style="single")' sample.yml
```
will output
```yaml
a:
  deeply:
    nested: 'newValue'
```

## Update multiple deeply nested properties
Given a sample.yml file of:
```yaml
a:
  deeply:
    nested: value
    other: thing
```
then
```bash
yq eval 'with(.a.deeply; .nested = "newValue" | .other= "newThing")' sample.yml
```
will output
```yaml
a:
  deeply:
    nested: newValue
    other: newThing
```

## Update array elements relatively
The second expression runs with each element of the array as it's contextual root. This allows you to make updates relative to the element.

Given a sample.yml file of:
```yaml
myArray:
  - a: apple
  - a: banana
```
then
```bash
yq eval 'with(.myArray[]; .b = .a + " yum")' sample.yml
```
will output
```yaml
myArray:
  - a: apple
    b: apple yum
  - a: banana
    b: banana yum
```

