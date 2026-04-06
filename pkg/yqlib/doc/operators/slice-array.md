# Slice Array or String

The slice operator works on both arrays and strings. Like the `jq` equivalent, `.[10:15]` will return a subarray (or substring) of length 5, starting from index 10 inclusive, up to index 15 exclusive. Negative numbers count backwards from the end of the array or string.

You may leave out the first or second number, which will refer to the start or end of the array or string respectively.

## Slicing arrays
Given a sample.yml file of:
```yaml
- cat
- dog
- frog
- cow
```
then
```bash
yq '.[1:3]' sample.yml
```
will output
```yaml
- dog
- frog
```

## Slicing arrays - without the first number
Starts from the start of the array

Given a sample.yml file of:
```yaml
- cat
- dog
- frog
- cow
```
then
```bash
yq '.[:2]' sample.yml
```
will output
```yaml
- cat
- dog
```

## Slicing arrays - without the second number
Finishes at the end of the array

Given a sample.yml file of:
```yaml
- cat
- dog
- frog
- cow
```
then
```bash
yq '.[2:]' sample.yml
```
will output
```yaml
- frog
- cow
```

## Slicing arrays - use negative numbers to count backwards from the end
Given a sample.yml file of:
```yaml
- cat
- dog
- frog
- cow
```
then
```bash
yq '.[1:-1]' sample.yml
```
will output
```yaml
- dog
- frog
```

## Inserting into the middle of an array
using an expression to find the index

Given a sample.yml file of:
```yaml
- cat
- dog
- frog
- cow
```
then
```bash
yq '(.[] | select(. == "dog") | key + 1) as $pos | .[0:($pos)] + ["rabbit"] + .[$pos:]' sample.yml
```
will output
```yaml
- cat
- dog
- rabbit
- frog
- cow
```

## Slicing strings
Given a sample.yml file of:
```yaml
country: Australia
```
then
```bash
yq '.country[0:5]' sample.yml
```
will output
```yaml
Austr
```

## Slicing strings - without the second number
Finishes at the end of the string

Given a sample.yml file of:
```yaml
country: Australia
```
then
```bash
yq '.country[5:]' sample.yml
```
will output
```yaml
alia
```

## Slicing strings - without the first number
Starts from the start of the string

Given a sample.yml file of:
```yaml
country: Australia
```
then
```bash
yq '.country[:5]' sample.yml
```
will output
```yaml
Austr
```

## Slicing strings - use negative numbers to count backwards from the end
Negative indices count from the end of the string

Given a sample.yml file of:
```yaml
country: Australia
```
then
```bash
yq '.country[-5:]' sample.yml
```
will output
```yaml
ralia
```

## Slicing strings - Unicode
Indices are rune-based, so multi-byte characters are handled correctly

Given a sample.yml file of:
```yaml
greeting: héllo
```
then
```bash
yq '.greeting[1:3]' sample.yml
```
will output
```yaml
él
```

