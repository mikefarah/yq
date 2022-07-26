# JSON

Encode and decode to and from JSON. Note that, unless you have multiple JSON documents in a single file, YAML is a _superset_ of JSON - so `yq` can read any json file without doing anything special.

If you do have mulitple JSON documents in a single file (e.g. NDJSON) then you will need to explicity use the json parser `-p=json`.

This means you don't need to 'convert' a JSON file to YAML - however if you want idiomatic YAML styling, then you can use the `-P/--prettyPrint` flag, see examples below.
