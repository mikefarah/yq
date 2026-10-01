# INI

Encode and decode to and from INI. Sections become maps, and keys without a section are added to the top level.

By default, both `=` and `:` are treated as key/value delimiters (matching the underlying INI parser's defaults). This can be controlled with the `--ini-key-value-delimiters` flag, which is useful when keys legitimately contain a colon, e.g. `this:key = value`.

See below for examples
