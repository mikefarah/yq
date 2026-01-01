# KYaml

Encode and decode to and from KYaml (a restricted subset of YAML that uses flow-style collections).

KYaml is useful when you want YAML data rendered in a compact, JSON-like form while still supporting YAML features like comments.

Notes:
- Strings are always double-quoted in KYaml output.
- Anchors and aliases are expanded (KYaml output does not emit them).
