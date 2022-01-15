# Select

Select is used to filter arrays and maps by a boolean expression.

## Select elements from array with regular expression
See more regular expression examples under the `string` operator docs.

Given a sample.yml file of:
```yaml
- this_0
- not_this
- nor_0_this
- thisTo_4
```
then
```bash
yq eval '.[] | select(test("[a-zA-Z]+_[0-9]$"))' sample.yml
```
will output
```yaml
this_0
thisTo_4
```

## Select elements from array using wildcard prefix
Given a sample.yml file of:
```yaml
- cat
- goat
- dog
```
then
```bash
yq eval '.[] | select(. == "*at")' sample.yml
```
will output
```yaml
cat
goat
```

## Select elements from array using wildcard suffix
Given a sample.yml file of:
```yaml
- go-kart
- goat
- dog
```
then
```bash
yq eval '.[] | select(. == "go*")' sample.yml
```
will output
```yaml
go-kart
goat
```

## Select elements from array using wildcard prefix and suffix
Given a sample.yml file of:
```yaml
- ago
- go
- meow
- going
```
then
```bash
yq eval '.[] | select(. == "*go*")' sample.yml
```
will output
```yaml
ago
go
going
```

## Select items from a map
Given a sample.yml file of:
```yaml
things: cat
bob: goat
horse: dog
```
then
```bash
yq eval '.[] | select(. == "cat" or test("og$"))' sample.yml
```
will output
```yaml
cat
dog
```

## Use select + with_entries to filter map keys
Given a sample.yml file of:
```yaml
name: bob
legs: 2
game: poker
```
then
```bash
yq eval 'with_entries(select(.key | test("ame$")))' sample.yml
```
will output
```yaml
name: bob
game: poker
```

## Select multiple items in a map and update
Note the brackets around the entire LHS.

Given a sample.yml file of:
```yaml
a:
  things: cat
  bob: goat
  horse: dog
```
then
```bash
yq eval '(.a.[] | select(. == "cat" or . == "goat")) |= "rabbit"' sample.yml
```
will output
```yaml
a:
  things: rabbit
  bob: rabbit
  horse: dog
```

