# Operators

In `yq` expressions are made up of operators and pipes. A context of nodes is passed through the expression and each operation takes the context as input and returns a new context as output. That output is piped in as input for the next operation in the expression. To begin with, the context is set to the first yaml document of the first yaml file (if processing in sequence using eval).

Lets look at a couple of examples.

## Example 1 - simple example

Given a document like:

```yaml
- [a]
- "cat"
```

with an expression:

```
.[] | length
```

`yq` will initially set the context as single node of the entire yaml document, an array of two elements.

```yaml
- [a]
- "cat"
```

This gets piped into the splat operator `.[]` which will split out the context into a collection of two nodes `[a]` and `"cat"`. Note that this is _not_ a yaml array.

The `length` operator take no arguments, and will simply return the length of _each_ matching node in the context. So for the context of `[a]` and `"cat"`, it will return a new context of `1` and `3`.

This being the last operation in the expression, the results will be printed out:

```
1
3
```

# Example 2 - operators with arguments.


The `=` operator takes two arguments, a `lhs` expression and `rhs` expression. It runs the 'matching' nodes context against the `lhs` expression to find the nodes to update, lets call it `lhsNodes`, and then runs the matching nodes against the `rhs` to find the new values, lets call that `rhsNodes`. It updates the `lhsNodes` values with the `rhsNodes` values and _returns the original matching nodes_. This is important, where length changed the matching nodes to be new nodes with the length values, `=` returns the original matching nodes, albeit with some of the nodes values updated. So `.a = 3` will still return the parent matching node, but with the matching child updated.

Please see the individual operator docs for more information and examples.

