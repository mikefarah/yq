
## Encode shell variables
Note that comments are dropped and values will be enclosed in single quotes as needed.

Given a sample.yml file of:
```yaml
# comment
name: Mike Wazowski
eyes:
  color: turquoise
  number: 1
friends:
  - James P. Sullivan
  - Celia Mae
```
then
```bash
yq -o=shell sample.yml
```
will output
```sh
name='Mike Wazowski'
eyes_color=turquoise
eyes_number=1
friends_0='James P. Sullivan'
friends_1='Celia Mae'
```

## Encode shell variables: illegal variable names as key.
Keys that would be illegal as variable keys are adapted.

Given a sample.yml file of:
```yaml
ascii_=_symbols: replaced with _
"ascii_	_controls": dropped (this example uses \t)
nonascii_א_characters: dropped
effrot_expeñded_tò_preserve_accented_latin_letters: moderate (via unicode NFKD)

```
then
```bash
yq -o=shell sample.yml
```
will output
```sh
ascii___symbols='replaced with _'
ascii__controls='dropped (this example uses \t)'
nonascii__characters=dropped
effrot_expended_to_preserve_accented_latin_letters='moderate (via unicode NFKD)'
```

## Encode shell variables: empty values, arrays and maps
Empty values are encoded to empty variables, but empty arrays and maps are skipped.

Given a sample.yml file of:
```yaml
empty:
  value:
  array: []
  map:   {}
```
then
```bash
yq -o=shell sample.yml
```
will output
```sh
empty_value=
```

## Encode shell variables: single quotes in values
Single quotes in values are encoded as '"'"' (close single quote, double-quoted single quote, open single quote).

Given a sample.yml file of:
```yaml
name: Miles O'Brien
```
then
```bash
yq -o=shell sample.yml
```
will output
```sh
name='Miles O'"'"'Brien'
```

