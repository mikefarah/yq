### Keys with dots
When specifying a key that has a dot use key lookup indicator.

```yaml
b:
  foo.bar: 7
```

```bash
yaml r sample.yaml 'b[foo.bar]'
```

```bash
yaml w sample.yaml 'b[foo.bar]' 9
```

Any valid yaml key can be specified as part of a key lookup.

Note that the path is in quotes to avoid the square brackets being interpreted by your shell.

### Keys (and values) with leading dashes
If a key or value has leading dashes, yq won't know that you are passing a value as opposed to a flag (and you will get a 'bad flag syntax' error).

To fix that, you will need to tell it to stop processing flags by adding '--' after the last flag like so:


```bash
yq n -t -- --key --value
```

Will result in

```
--key: --value
```