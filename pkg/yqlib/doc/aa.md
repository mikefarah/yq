# Operators

In `yq` expressions are made up of operators. Operators have 0-2 arguments and run against the current 'matching' nodes in the expression tree.

Lets look at a couple of examples.

The `length` operator take no arguments, and will simply return the length of _each_ matching node. So if there were 2 nodes, one string and one array, length will update the 'matching' nodes context to be two new numeric scalar nodes representing the lengths of the orignal 'matching' nodes.

The `=` operator takes two arguments, a `lhs` expression and `rhs` expression. It runs the 'matching' nodes context against the `lhs` expression to find the nodes to update, lets call it `lhsNodes`, and then runs the matching nodes against the `rhs` to find the new values, lets call that `rhsNodes`. It updates the `lhsNodes` values with the `rhsNodes` values and _returns the original matching nodes_. This is important, where length changed the matching nodes to be new nodes with the length values, `=` returns the original matching nodes, albeit with some of the nodes values updated. So `.a = 3` will still return the parent matching node, but with the matching child updated.

Please see the individual operator docs for more information and examples.

