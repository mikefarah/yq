# INI

Encode and decode to and from INI. Sections become maps, and keys without a section are added to the top level.

By default, both `=` and `:` are treated as key/value delimiters (matching the underlying INI parser's defaults). This can be controlled with the `--ini-key-value-delimiters` flag, which is useful when keys legitimately contain a colon, e.g. `this:key = value`.

See below for examples

## Parse INI: simple
Given a sample.ini file of:
```ini
[section]
key = value

```
then
```bash
yq -p=ini sample.ini
```
will output
```yaml
section:
  key: value
```

## Encode INI: simple
Given a sample.yml file of:
```yaml
section: {key: value}
```
then
```bash
yq -o=ini '.' sample.yml
```
will output
```ini
[section]
key = value
```

## Roundtrip INI: simple
Given a sample.ini file of:
```ini
[section]
key = value

```
then
```bash
yq -p=ini -o=ini sample.ini
```
will output
```ini
[section]
key = value
```

## bad ini
Given a sample.ini file of:
```ini
[section\nkey = value
```
then an error is expected:
```
bad file 'sample.yml': failed to parse INI content: unclosed section: [section\nkey = value
```

## Parse INI: key with colon
By default, the key/value delimiters are "=:", so ':' is treated the same as '=' and the key is split at the first delimiter found. See the next example for how to avoid this using `--ini-key-value-delimiters`.

Given a sample.ini file of:
```ini
[this:section]
this:line = should really work

```
then
```bash
yq -p=ini sample.ini
```
will output
```yaml
this:section:
  this: line = should really work
```

## Parse INI: key with colon, using --ini-key-value-delimiters
Keys containing a colon (e.g. Mercurial config files) are parsed correctly when ':' is removed from the key/value delimiters.

Given a sample.ini file of:
```ini
[this:section]
this:line = should really work

```
then
```bash
yq -p=ini --ini-key-value-delimiters='=' sample.ini
```
will output
```yaml
this:section:
  this:line: should really work
```

instead of
```yaml
this:section:
  this: line = should really work
```

