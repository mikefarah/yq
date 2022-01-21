# JSON

Encode and decode to and from JSON. Note that YAML is a _superset_ of JSON - so `yq` can read any json file without doing anything special.

This means you don't need to 'convert' a JSON file to YAML - however if you want idiomatic YAML styling, then you can use the `-P/--prettyPrint` flag, see examples below.
