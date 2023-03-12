# Eval

Use `eval` to dynamically process an expression - for instance from an environment variable.

`eval` takes a single argument, and evaluates that as a `yq` expression. Any valid expression can be used, be it a path `.a.b.c | select(. == "cat")`, or an update `.a.b.c = "gogo"`.

Tip: This can be useful way parameterise complex scripts.
