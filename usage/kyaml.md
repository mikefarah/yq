# KYaml

Encode and decode to and from KYaml (a restricted subset of YAML that uses flow-style collections).

KYaml is useful when you want YAML data rendered in a compact, JSON-like form while still supporting YAML features like comments.

Notes:
- Strings are always double-quoted in KYaml output.
- Anchors and aliases are expanded (KYaml output does not emit them).

## Encode kyaml: plain string scalar
Strings are always double-quoted in KYaml output.

Given a sample.yml file of:
```yaml
cat

```
then
```bash
yq -o=kyaml '.' sample.yml
```
will output
```yaml
"cat"
```

## encode flow mapping and sequence
Given a sample.yml file of:
```yaml
a: b
c:
  - d

```
then
```bash
yq -o=kyaml '.' sample.yml
```
will output
```yaml
{
  a: "b",
  c: [
    "d",
  ],
}
```

## encode non-string scalars
Given a sample.yml file of:
```yaml
a: 12
b: true
c: null
d: "true"

```
then
```bash
yq -o=kyaml '.' sample.yml
```
will output
```yaml
{
  a: 12,
  b: true,
  c: null,
  d: "true",
}
```

## quote non-identifier keys
Given a sample.yml file of:
```yaml
"1a": b
"has space": c

```
then
```bash
yq -o=kyaml '.' sample.yml
```
will output
```yaml
{
  "1a": "b",
  "has space": "c",
}
```

## escape quoted strings
Given a sample.yml file of:
```yaml
a: "line1\nline2\t\"q\""

```
then
```bash
yq -o=kyaml '.' sample.yml
```
will output
```yaml
{
  a: "line1\nline2\t\"q\"",
}
```

## preserve comments when encoding
Given a sample.yml file of:
```yaml
# leading
a: 1 # a line
# head b
b: 2
c:
  # head d
  - d # d line
  - e
# trailing

```
then
```bash
yq -o=kyaml '.' sample.yml
```
will output
```yaml
# leading
{
  a: 1, # a line
  # head b
  b: 2,
  c: [
    # head d
    "d", # d line
    "e",
  ],
  # trailing
}
```

## Encode kyaml: anchors and aliases
KYaml output does not support anchors/aliases; they are expanded to concrete values.

Given a sample.yml file of:
```yaml
base: &base
  a: b
copy: *base

```
then
```bash
yq -o=kyaml '.' sample.yml
```
will output
```yaml
{
  base: {
    a: "b",
  },
  copy: {
    a: "b",
  },
}
```

## Encode kyaml: yaml to kyaml shows formatting differences
KYaml uses flow-style collections (braces/brackets) and explicit commas.

Given a sample.yml file of:
```yaml
person:
  name: John
  pets:
    - cat
    - dog

```
then
```bash
yq -o=kyaml '.' sample.yml
```
will output
```yaml
{
  person: {
    name: "John",
    pets: [
      "cat",
      "dog",
    ],
  },
}
```

## Encode kyaml: nested lists of objects
Lists and objects can be nested arbitrarily; KYaml always uses flow-style collections.

Given a sample.yml file of:
```yaml
- name: a
  items:
    - id: 1
      tags:
        - k: x
          v: y
        - k: x2
          v: y2
    - id: 2
      tags:
        - k: z
          v: w

```
then
```bash
yq -o=kyaml '.' sample.yml
```
will output
```yaml
[
  {
    name: "a",
    items: [
      {
        id: 1,
        tags: [
          {
            k: "x",
            v: "y",
          },
          {
            k: "x2",
            v: "y2",
          },
        ],
      },
      {
        id: 2,
        tags: [
          {
            k: "z",
            v: "w",
          },
        ],
      },
    ],
  },
]
```

