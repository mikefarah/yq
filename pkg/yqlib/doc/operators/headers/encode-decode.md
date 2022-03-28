# Encoder / Decoder

Encode operators will take the piped in object structure and encode it as a string in the desired format. The decode operators do the opposite, they take a formatted string and decode it into the relevant object structure.

Note that you can optionally pass an indent value to the encode functions (see below).

These operators are useful to process yaml documents that have stringified embeded yaml/json/props in them.


| Format | Decode (from string) | Encode (to string) |
| --- | -- | --|
| Yaml | from_yaml | to_yaml(i)/@yaml |
| JSON | from_json | to_json(i)/@json |
| Properties | from_props  | to_props/@props |
| CSV |  | to_csv/@csv |
| TSV |  | to_tsv/@tsv |
| XML | from_xml | to_xml(i)/@xml |
| Base64 | @base64d | @base64 |


CSV and TSV format both accept either a single array or scalars (representing a single row), or an array of array of scalars (representing multiple rows). 

XML uses the `--xml-attribute-prefix` and `xml-content-name` flags to identify attributes and content fields.


Base64 assumes [rfc4648](https://rfc-editor.org/rfc/rfc4648.html) encoding. Encoding and decoding both assume that the content is a string.
