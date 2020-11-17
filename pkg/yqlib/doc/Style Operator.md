
## Examples
### Example 0
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq eval '.a style="single"' sample.yml
```
will output
```yaml
{a: 'cat'}
```

### Set style using a path
Given a sample.yml file of:
```yaml
a: cat
b: double
```
then
```bash
yq eval '.a style=.b' sample.yml
```
will output
```yaml
{a: "cat", b: double}
```

### Example 2
Given a sample.yml file of:
```yaml
a: cat
b: dog
```
then
```bash
yq eval '.. style=""' sample.yml
```
will output
```yaml
a: cat
b: dog
```

### Example 3
Given a sample.yml file of:
```yaml
a: cat
b: thing
```
then
```bash
yq eval '.. | style' sample.yml
```
will output
```yaml
flow
double
single
```

### Example 4
Given a sample.yml file of:
```yaml
a: cat
```
then
```bash
yq eval '.. | style' sample.yml
```
will output
```yaml


```

