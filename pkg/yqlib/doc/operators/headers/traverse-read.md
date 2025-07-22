# Traverse (Read)

This is the simplest (and perhaps most used) operator. It is used to navigate deeply into yaml structures.


## NOTE --yaml-fix-merge-anchor-to-spec flag
`yq` doesn't merge anchors `<<:` to spec, in some circumstances it incorrectly overrides existing keys when the spec documents not to do that.

To minimise disruption while still fixing the issue, a flag has been added to toggle this behaviour. This will first default to false; and log warnings to users. Then it will default to true (and still allow users to specify false if needed)

See examples of the flag differences below.

