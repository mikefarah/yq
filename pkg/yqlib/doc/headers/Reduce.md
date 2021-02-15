Reduce is a powerful way to process a collection of data into a new form.

## yq vs jq syntax
Reduce syntax in `yq` is a little different from `jq` - as `yq` (currently) isn't as sophisticated as `jq` and its only supports infix notation (e.g. a + b, the operator is in the middle of the two parameters) - where as `jq` uses a mix of infix notation with _prefix_ notation (e.g. `reduce a b` is like writing `+ a b`).

To that end, the reduce operator is called `ireduce` for backwards compatability if a prefix version of `reduce` is ever added.

