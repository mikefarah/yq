# Comment Operators

Use these comment operators to set or retrieve comments. Note that line comments on maps/arrays are actually set on the _key_ node as opposed to the _value_ (map/array). See below for examples.

Like the `=` and `|=` assign operators, the same syntax applies when updating comments:

### plain form: `=`
This will set the LHS nodes' comments equal to the expression on the RHS. The RHS is run against the matching nodes in the pipeline

### relative form: `|=` 
This is similar to the plain form, but it evaluates the RHS with _each matching LHS node as context_. This is useful if you want to set the comments as a relative expression of the node, for instance its value or path.
