# Load

The load operators allows you to load in content from another file.

Note that you can use string operators like `+` and `sub` to modify the value in the yaml file to a path that exists in your system.

You can load files of the following supported types:

|Format | Load Operator |
| --- | --- |
| Yaml | load |
| XML | load_xml |
| Properties | load_props |
| Plain String | load_str |
| Base64 | load_base64 |

Note that load_base64 only works for base64 encoded utf-8 strings.

## Samples files for tests:

### yaml

`../../examples/thing.yml`:

```yaml
a: apple is included
b: cool
```

### xml
`small.xml`:

```xml
<this>is some xml</this>
```

### properties
`small.properties`:

```properties
this.is = a properties file
```

### base64
`base64.txt`:
```
bXkgc2VjcmV0IGNoaWxsaSByZWNpcGUgaXMuLi4u
```

## Disabling file operators
If required, you can use the `--security-disable-file-ops` to disable file operations.

