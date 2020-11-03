# Equal Operator
## Examples
### Example 0
sample.yml:
```yaml
[cat,goat,dog]
```
Expression
```bash
yq '.[] | (. == "*at")' < sample.yml
```
Result
```yaml
true
true
false
```
### Example 1
sample.yml:
```yaml
[3, 4, 5]
```
Expression
```bash
yq '.[] | (. == 4)' < sample.yml
```
Result
```yaml
false
true
false
```
### Example 2
sample.yml:
```yaml
a: { cat: {b: apple, c: whatever}, pat: {b: banana} }
```
Expression
```bash
yq '.a | (.[].b == "apple")' < sample.yml
```
Result
```yaml
true
false
```
### Example 3
Expression
```bash
yq 'null == null' < sample.yml
```
Result
```yaml
true
```
### Example 4
Expression
```bash
yq 'null == ~' < sample.yml
```
Result
```yaml
true
```
