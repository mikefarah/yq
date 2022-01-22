# XML

Encode and decode to and from XML. Whitespace is not conserved for round trips - but the order of the fields are.

Consecutive xml nodes with the same name are assumed to be arrays.

XML content data and attributes are created as fields. This can be controlled by the `'--xml-attribute-prefix` and `--xml-content-name` flags - see below for examples.
