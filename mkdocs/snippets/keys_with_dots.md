### Keys with dots
When specifying a key that has a dot use key lookup indicator.

```yaml
b:
  foo.bar: 7
```

```bash
yaml r sample.yaml b[foo.bar]
```

```bash
yaml w sample.yaml b[foo.bar] 9
```

Any valid yaml key can be specified as part of a key lookup.
