
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

