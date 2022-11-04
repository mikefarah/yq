# Sort

Sorts an array. Use `sort` to sort an array as is, or `sort_by(exp)` to sort by a particular expression (e.g. subfield).

To sort by descending order, pipe the results through the `reverse` operator after sorting.

Note that at this stage, `yq` only sorts scalar fields.


## Sort by string field
Given a sample.yml file of:
```yaml
- a: banana
- a: cat
- a: apple
```
then
```bash
yq 'sort_by(.a)' sample.yml
```
will output
```yaml
- a: apple
- a: banana
- a: cat
```

## Sort descending by string field
Use sort with reverse to sort in descending order.

Given a sample.yml file of:
```yaml
- a: banana
- a: cat
- a: apple
```
then
```bash
yq 'sort_by(.a) | reverse' sample.yml
```
will output
```yaml
- a: cat
- a: banana
- a: apple
```

## Sort array in place
Given a sample.yml file of:
```yaml
cool:
  - a: banana
  - a: cat
  - a: apple
```
then
```bash
yq '.cool |= sort_by(.a)' sample.yml
```
will output
```yaml
cool:
  - a: apple
  - a: banana
  - a: cat
```

## Sort array of objects by key
Note that you can give sort_by complex expressions, not just paths

Given a sample.yml file of:
```yaml
cool:
  - b: banana
  - a: banana
  - c: banana
```
then
```bash
yq '.cool |= sort_by(keys | .[0])' sample.yml
```
will output
```yaml
cool:
  - a: banana
  - b: banana
  - c: banana
```

## Sort is stable
Note the order of the elements in unchanged when equal in sorting.

Given a sample.yml file of:
```yaml
- a: banana
  b: 1
- a: banana
  b: 2
- a: banana
  b: 3
- a: banana
  b: 4
```
then
```bash
yq 'sort_by(.a)' sample.yml
```
will output
```yaml
- a: banana
  b: 1
- a: banana
  b: 2
- a: banana
  b: 3
- a: banana
  b: 4
```

## Sort by numeric field
Given a sample.yml file of:
```yaml
- a: 10
- a: 100
- a: 1
```
then
```bash
yq 'sort_by(.a)' sample.yml
```
will output
```yaml
- a: 1
- a: 10
- a: 100
```

## Sort by custom date field
Given a sample.yml file of:
```yaml
- a: 12-Jun-2011
- a: 23-Dec-2010
- a: 10-Aug-2011
```
then
```bash
yq 'with_dtf("02-Jan-2006"; sort_by(.a))' sample.yml
```
will output
```yaml
- a: 23-Dec-2010
- a: 12-Jun-2011
- a: 10-Aug-2011
```

## Sort, nulls come first
Given a sample.yml file of:
```yaml
- 8
- 3
- null
- 6
- true
- false
- cat
```
then
```bash
yq 'sort' sample.yml
```
will output
```yaml
- null
- false
- true
- 3
- 6
- 8
- cat
```

