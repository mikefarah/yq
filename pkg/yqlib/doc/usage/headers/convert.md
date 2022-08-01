# JSON

Encode and decode to and from JSON. Supports multiple JSON documents in a single file (e.g. NDJSON).

Note that YAML is a superset of (single document) JSON - so you don't have to use the JSON parser to read JSON when there is only one JSON document in the input. You will probably want to pretty print the result in this case, to get idiomatic YAML styling.

