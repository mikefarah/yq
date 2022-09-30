# Comment Operators

Use these comment operators to set or retrieve comments. Note that line comments on maps/arrays are actually set on the _key_ node as opposed to the _value_ (map/array). See below for examples.

Like the `=` and `|=` assign operators, the same syntax applies when updating comments:

### plain form: `=`
This will assign the LHS nodes comments to the expression on the RHS. The RHS is run against the matching nodes in the pipeline

### relative form: `|=` 
Similar to the plain form, however the RHS evaluates against each matching LHS node! This is useful if you want to set the comments as a relative expression of the node, for instance its value or path.
