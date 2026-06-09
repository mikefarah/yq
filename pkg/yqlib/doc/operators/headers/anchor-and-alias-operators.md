# Anchor and Alias Operators

Use the `alias` and `anchor` operators to read and write yaml aliases and anchors. The `explode` operator normalises a yaml file (dereference (or expands) aliases and remove anchor names).

`yq` supports merge keys (like `<<: *blah`) from YAML 1.1. These are no longer part of the YAML 1.2 standard, but remain common in practice. Plain `<<:` keys are recognised as merge keys and round-trip as `<<:` without an explicit `!!merge` tag. When the source uses an explicit `!!merge` tag, that is preserved on output. Internally, when `yq` synthesises a `<<` map key (for example during merge operations), it tags the key as `!!merge` rather than `!!str`.


## NOTE --yaml-fix-merge-anchor-to-spec flag
`yq` doesn't merge anchors `<<:` to spec, in some circumstances it incorrectly overrides existing keys when the spec documents not to do that.

To minimise disruption while still fixing the issue, a flag has been added to toggle this behaviour. This will first default to false; and log warnings to users. Then it will default to true (and still allow users to specify false if needed).

This flag also enables advanced merging, like inline maps, as well as fixes to ensure when exploding a particular path, neighbours are not affect ed.

Long story short, you should be setting this flag to true.

See examples of the flag differences below, where LEGACY is with the flag off; and FIXED is with the flag on.

