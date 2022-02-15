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

Lets say there is a file `../../examples/thing.yml`:

```yaml
a: apple is included
b: cool
```
and a file `small.xml`:

```xml
<this>is some xml</this>
```

and `small.properties`:

```properties
this.is = a properties file
```
