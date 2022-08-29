# Assign (Update)

This operator is used to update node values. It can be used in either the:

### plain form: `=`
Which will assign the LHS node values to the RHS node values. The RHS expression is run against the matching nodes in the pipeline.

### relative form: `|=`
This will do a similar thing to the plain form, however, the RHS expression is run against _the LHS nodes_. This is useful for updating values based on old values, e.g. increment.


### Flags
- `c` clobber custom tags
