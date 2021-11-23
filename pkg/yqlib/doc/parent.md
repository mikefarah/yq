
## Simple example
Given a sample.yml file of:
```yaml
a:
  nested: cat
```
then
```bash
yq eval '.a.nested | parent' sample.yml
```
will output
```yaml
nested: cat
```

## Show parent
Given a sample.yml file of:
```yaml
a:
  fruit: apple
b:
  fruit: banana
```
then
```bash
yq eval '.. | select(. == "banana") | parent' sample.yml
```
will output
```yaml
fruit: banana
```

## No parent
Given a sample.yml file of:
```yaml
{}
```
then
```bash
yq eval 'parent' sample.yml
```
will output
```yaml
```

