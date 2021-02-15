Reduce is a powerful way to process a collection of data into a new form.

```
<exp> as $<name> ireduce (<init>; <block>)
```

e.g.

```
.[] as $item ireduce (0; . + $item)
```

On the LHS we are configuring the collection of items that will be reduced `<exp>` as well as what each element will be called `$<name>`. Note that the array has been splatted into its individual elements.

On the RHS there is `<init>`, the starting value of the accumulator and `<block>`, the expression that will update the accumulator for each element in the collection. 

Note that within the block expression, `.` will evaluate to the current value of the accumulator. This effectively means that within the `reduce` block you can no longer access data other than elements of the array set as `$<name>`. For simple things, this is probably fine, but often you will need to refer to other data elements.

This can be done by setting a variable using `as` and piping that into the `reduce` operation, or you can simply refer to `$context` which is exactly that, automatically set for you for convenience. See examples below.

## yq vs jq syntax
Reduce syntax in `yq` is a little different from `jq` - as `yq` (currently) isn't as sophisticated as `jq` and its only supports infix notation (e.g. a + b, where the operator is in the middle of the two parameters) - where as `jq` uses a mix of infix notation with _prefix_ notation (e.g. `reduce a b` is like writing `+ a b`).

To that end, the reduce operator is called `ireduce` for backwards compatability if a prefix version of `reduce` is ever added.