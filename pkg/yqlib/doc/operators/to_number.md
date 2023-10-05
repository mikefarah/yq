# To Number
Parses the input as a number. yq will try to parse values as an int first, failing that it will try float. Values that already ints or floats will be left alone.

## Converts strings to numbers
Given a sample.yml file of:
```yaml
- "3"
- "3.1"
- "-1e3"
```
then
```bash
yq '.[] | to_number' sample.yml
```
will output
```yaml
3
3.1
-1e3
```

## Doesn't change numbers
Given a sample.yml file of:
```yaml
- 3
- 3.1
- -1e3
```
then
```bash
yq '.[] | to_number' sample.yml
```
will output
```yaml
3
3.1
-1e3
```

## Cannot convert null
Running
```bash
yq --null-input '.a.b | to_number'
```
will output
```bash
Error: cannot convert node value [null] at path a.b of tag !!null to number
```

