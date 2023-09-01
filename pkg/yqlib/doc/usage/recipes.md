# Recipes

These examples are intended to show how you can use multiple operators together so you get an idea of how you can perform complex data manipulation.

Please see the details [operator docs](https://mikefarah.gitbook.io/yq/operators) for details on each individual operator.

## Find items in an array
We have an array and we want to find the elements with a particular name.

Given a sample.yml file of:
```yaml
- name: Foo
  numBuckets: 0
- name: Bar
  numBuckets: 0
```
then
```bash
yq '.[] | select(.name == "Foo")' sample.yml
```
will output
```yaml
name: Foo
numBuckets: 0
```

### Explanation:
- `.[]` splats the array, and puts all the items in the context.
- These items are then piped (`|`) into `select(.name == "Foo")` which will select all the nodes that have a name property set to 'Foo'.
- See the [select](https://mikefarah.gitbook.io/yq/operators/select) operator for more information.

## Find and update items in an array
We have an array and we want to _update_ the elements with a particular name.

Given a sample.yml file of:
```yaml
- name: Foo
  numBuckets: 0
- name: Bar
  numBuckets: 0
```
then
```bash
yq '(.[] | select(.name == "Foo") | .numBuckets) |= . + 1' sample.yml
```
will output
```yaml
- name: Foo
  numBuckets: 1
- name: Bar
  numBuckets: 0
```

### Explanation:
- Following from the example above`.[]` splats the array, selects filters the items.
- We then pipe (`|`) that into `.numBuckets`, which will select that field from all the matching items
- Splat, select and the field are all in brackets, that whole expression is passed to the `|=` operator as the left hand side expression, with `. + 1` as the right hand side expression.
- `|=` is the operator that updates fields relative to their own value, which is referenced as dot (`.`).
- The expression `. + 1` increments the numBuckets counter.
- See the [assign](https://mikefarah.gitbook.io/yq/operators/assign-update) and [add](https://mikefarah.gitbook.io/yq/operators/add) operators for more information.

## Multiple or complex updates to items in an array
We have an array and we want to _update_ the elements with a particular name in reference to its type.

Given a sample.yml file of:
```yaml
myArray:
  - name: Foo
    type: cat
  - name: Bar
    type: dog
```
then
```bash
yq 'with(.myArray[]; .name = .name + " - " + .type)' sample.yml
```
will output
```yaml
myArray:
  - name: Foo - cat
    type: cat
  - name: Bar - dog
    type: dog
```

### Explanation:
- The with operator will effectively loop through each given item in the first given expression, and run the second expression against it.
- `.myArray[]` splats the array in `myArray`. So `with` will run against each item in that array
- `.name = .name + " - " + .type` this expression is run against every item, updating the name to be a concatenation of the original name as well as the type.
- See the [with](https://mikefarah.gitbook.io/yq/operators/with) operator for more information and examples.

## Sort an array by a field
Given a sample.yml file of:
```yaml
myArray:
  - name: Foo
    numBuckets: 1
  - name: Bar
    numBuckets: 0
```
then
```bash
yq '.myArray |= sort_by(.numBuckets)' sample.yml
```
will output
```yaml
myArray:
  - name: Bar
    numBuckets: 0
  - name: Foo
    numBuckets: 1
```

### Explanation:
- We want to resort `.myArray`.
- `sort_by` works by piping an array into it, and it pipes out a sorted array.
- So, we use `|=` to update `.myArray`. This is the same as doing `.myArray = (.myArray | sort_by(.numBuckets))`

## Filter, flatten, sort and unique
Lets

Given a sample.yml file of:
```yaml
- type: foo
  names:
    - Fred
    - Catherine
- type: bar
  names:
    - Zelda
- type: foo
  names: Fred
- type: foo
  names: Ava
```
then
```bash
yq '[.[] | select(.type == "foo") | .names] | flatten | sort | unique' sample.yml
```
will output
```yaml
- Ava
- Catherine
- Fred
```

### Explanation:
- `.[] | select(.type == "foo") | .names` will select the array elements of type "foo"
- Splat `.[]` will unwrap the array and match all the items. We need to do this so we can work on the child items, for instance, filter items out using the `select` operator.
- But we still want the final results back into an array. So after we're doing working on the children, we wrap everything back into an array using square brackets around the expression. `[.[] | select(.type == "foo") | .names]`
- Now have have an array of all the 'names' values. Which includes arrays of strings as well as strings on their own.
- Pipe `|` this array through `flatten`. This will flatten nested arrays. So now we have a flat list of all the name value strings
- Next we pipe `|` that through `sort` and then `unique` to get a sorted, unique list of the names!
- See the [flatten](https://mikefarah.gitbook.io/yq/operators/flatten), [sort](https://mikefarah.gitbook.io/yq/operators/sort) and [unique](https://mikefarah.gitbook.io/yq/operators/unique) for more information and examples.

