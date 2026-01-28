# Slice/Splice Array

The slice array operator takes an array as input and returns a subarray. Like the `jq` equivalent, `.[10:15]` will return an array of length 5, starting from index 10 inclusive, up to index 15 exclusive. Negative numbers count backwards from the end of the array.

You may leave out the first or second number, which will refer to the start or end of the array respectively.

{% hint style="warning" %}

_Note_ that slicing only works on arrays, not strings. Using slice syntax on a string (e.g., `.text | .[0:5]`) will return an empty result instead of a substring.

For string slicing, use split/join or regex:

```bash
# Split into array, slice, rejoin
yq '.text | split("") | .[0:5] | join("")'

# Or use regex substitution
yq '.text | sub("^(.{5}).*"; "${1}")'
```

{% endhint %}

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

