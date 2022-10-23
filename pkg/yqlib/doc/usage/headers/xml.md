# XML

Encode and decode to and from XML. Whitespace is not conserved for round trips - but the order of the fields are.

Consecutive xml nodes with the same name are assumed to be arrays.

XML content data, attributes processing instructions and directives are all created as plain fields. 

This can be controlled by:

| Flag | Default |Sample XML | 
| -- | -- |  -- |
 | `--xml-attribute-prefix` | `+` (changing to `+@` soon) | Legs in ```<cat legs="4"/>``` |  
 |  `--xml-content-name` | `+content` | Meow in ```<cat>Meow <fur>true</true></cat>``` |
 | `--xml-directive-name` | `+directive` | ```<!DOCTYPE config system "blah">``` |
 | `--xml-proc-inst-prefix` | `+p_` |  ```<?xml version="1"?>``` |


{% hint style="warning" %}
Default Attribute Prefix will be changing in v4.30!
In order to avoid name conflicts (e.g. having an attribute named "content" will create a field that clashes with the default content name of "+content") the attribute prefix will be changing to "+@".

This will affect users that have not set their own prefix and are not roundtripping XML changes.

{% endhint %}

## Encoder / Decoder flag options

In addition to the above flags, there are the following xml encoder/decoder options controlled by flags:

| Flag | Default | Description |
| -- | -- | -- |
| `--xml-strict-mode` | false | Strict mode enforces the requirements of the XML specification. When switched off the parser allows input containing common mistakes. See [the Golang xml decoder ](https://pkg.go.dev/encoding/xml#Decoder) for more details.| 
| `--xml-keep-namespace` | true | Keeps the namespace of attributes |
| `--xml-raw-token` | true |  Does not verify that start and end elements match and does not translate name space prefixes to their corresponding URLs. |
| `--xml-skip-proc-inst` | false | Skips over processing instructions, e.g. `<?xml version="1"?>` |
| `--xml-skip-directives` | false | Skips over directives, e.g. ```<!DOCTYPE config system "blah">``` |


See below for examples
