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

## Deeply prune a tree
Say we are only interested in child1 and child2, and want to filter everything else out.

Given a sample.yml file of:
```yaml
parentA:
  - bob
parentB:
  child1: i am child1
  child3: hiya
parentC:
  childX: cool
  child2: me child2
```
then
```bash
yq '(
  .. | # recurse through all the nodes
  select(has("child1") or has("child2")) | # match parents that have either child1 or child2
  (.child1, .child2) | # select those children
  select(.) # filter out nulls
) as $i ireduce({};  # using that set of nodes, create a new result map
  setpath($i | path; $i) # and put in each node, using its original path
)' sample.yml
```
will output
```yaml
parentB:
  child1: i am child1
parentC:
  child2: me child2
```

### Explanation:
- Find all the matching child1 and child2 nodes
- Using ireduce, create a new map using just those nodes
- Set each node into the new map using its original path

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
Lets find the unique set of names from the document.

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

## Export as environment variables (script), or any custom format
Given a yaml document, lets output a script that will configure environment variables with that data. This same approach can be used for exporting into custom formats.

Given a sample.yml file of:
```yaml
var0: string0
var1: string1
fruit:
  - apple
  - banana
  - peach
```
then
```bash
yq '.[] |(
	( select(kind == "scalar") | key + "='\''" + . + "'\''"),
	( select(kind == "seq") | key + "=(" + (map("'\''" + . + "'\''") | join(",")) + ")")
)' sample.yml
```
will output
```yaml
var0='string0'
var1='string1'
fruit=('apple','banana','peach')
```

### Explanation:
- `.[]` matches all top level elements
- We need a string expression for each of the different types that will produce the bash syntax, we'll use the union operator, to join them together
- Scalars, we just need the key and quoted value: `( select(kind == "scalar") | key + "='" + . + "'")`
- Sequences (or arrays) are trickier, we need to quote each value and `join` them with `,`: `map("'" + . + "'") | join(",")`

## Custom format with nested data
Like the previous example, but lets handle nested data structures. In this custom example, we're going to join the property paths with _. The important thing to keep in mind is that our expression is not recursive (despite the data structure being so). Instead we match _all_ elements on the tree and operate on them.

Given a sample.yml file of:
```yaml
simple: string0
simpleArray:
  - apple
  - banana
  - peach
deep:
  property: value
  array:
    - cat
```
then
```bash
yq '.. |(
	( select(kind == "scalar" and parent | kind != "seq") | (path | join("_")) + "='\''" + . + "'\''"),
	( select(kind == "seq") | (path | join("_")) + "=(" + (map("'\''" + . + "'\''") | join(",")) + ")")
)' sample.yml
```
will output
```yaml
simple='string0'
deep_property='value'
simpleArray=('apple','banana','peach')
deep_array=('cat')
```

### Explanation:
- You'll need to understand how the previous example works to understand this extension.
- `..` matches _all_ elements, instead of `.[]` from the previous example that just matches top level elements.
- Like before, we need a string expression for each of the different types that will produce the bash syntax, we'll use the union operator, to join them together
- This time, however, our expression matches every node in the data structure.
- We only want to print scalars that are not in arrays (because we handle the separately), so well add `and parent | kind != "seq"` to the select operator expression for scalars
- We don't just want the key any more, we want the full path. So instead of `key` we have `path | join("_")`
- The expression for sequences follows the same logic

