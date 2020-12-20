This is a boolean operator that will return ```true``` if the LHS is equal to the RHS and ``false`` otherwise.

```
.a == .b
```

It is most often used with the select operator to find particular nodes:

```
select(.a == .b)
```


## Match string
Given a sample.yml file of:
```yaml
- cat
- goat
- dog
```
then
```bash
yq eval '.[] | (. == "*at")' sample.yml
```
will output
```yaml
true
true
false
```

## Match number
Given a sample.yml file of:
```yaml
- 3
- 4
- 5
```
then
```bash
yq eval '.[] | (. == 4)' sample.yml
```
will output
```yaml
false
true
false
```

## Match nulls
Running
```bash
yq eval --null-input 'null == ~'
```
will output
```yaml
true
```

