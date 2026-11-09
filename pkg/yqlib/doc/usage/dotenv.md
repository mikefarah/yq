
## Decode a dotenv file
Comments and blank lines are dropped and an `export` prefix is accepted. The value is a string, so `1` stays quoted in YAML output.

Given a sample.env file of:
```env
# a comment
NAME=Mike Wazowski
export ROLE=scarer

EYES=1 # inline comment

```
then
```bash
yq -p=dotenv '.' sample.env
```
will output
```yaml
NAME: Mike Wazowski
ROLE: scarer
EYES: "1"
```

## Decode quoted values
`\n` is the only escape sequence, so `\t` stays two characters. A single-quoted value is taken literally and a quote ends a double-quoted value.

Given a sample.env file of:
```env
NEWLINE="line one\nline two"
TAB="a\tb"
SINGLE='line one\nline two'
HASH="a # b"
TRUNCATED="say \"hi\""

```
then
```bash
yq -p=dotenv '.' sample.env
```
will output
```yaml
NEWLINE: |-
  line one
  line two
TAB: a\tb
SINGLE: line one\nline two
HASH: 'a # b'
TRUNCATED: say \
```

## Decode a multi-line value
A quoted value may span lines, which is how a PEM key is written in a .env file.

Given a sample.env file of:
```env
KEY="-----BEGIN-----
abc
-----END-----"
NEXT=1

```
then
```bash
yq -p=dotenv '.' sample.env
```
will output
```yaml
KEY: |-
  -----BEGIN-----
  abc
  -----END-----
NEXT: "1"
```

## Decode the awkward lines
CRLF endings, a byte order mark, an `#` inside an unquoted value (which starts a comment, as in Node), an `=` in the value, and an empty value.

Given a sample.env file of:
```env
﻿WINDOWS=crlf
HASH=a#b
EQUALS=a=b
EMPTY=

```
then
```bash
yq -p=dotenv '.' sample.env
```
will output
```yaml
WINDOWS: crlf
HASH: a
EQUALS: a=b
EMPTY: ""
```

## Encode a dotenv file
Nested keys are joined with an underscore and sequences are indexed, as an environment variable cannot be nested.

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
yq -o=dotenv '.' sample.yml
```
will output
```env
name='Mike Wazowski'
eyes_color=turquoise
eyes_number=1
friends_0='James P. Sullivan'
friends_1='Celia Mae'
```

## Encode values that need quoting
Single quotes are used where every reader takes them literally; double quotes only for a value that spans lines, where `\n` is the escape.

Given a sample.yml file of:
```yaml
plain: /usr/local/bin
spaced: two words
hash: "a # b"
dollar: $HOME
backslash: "C:\\temp\\x"
quoted: 'say "hi"'
apostrophe: "it's fine"
multiline: |-
  one
  two
empty: ""

```
then
```bash
yq -o=dotenv '.' sample.yml
```
will output
```env
plain=/usr/local/bin
spaced='two words'
hash='a # b'
dollar='$HOME'
backslash='C:\temp\x'
quoted='say "hi"'
apostrophe="it's fine"
multiline="one\ntwo"
empty=""
```

## Round trip a dotenv file
Decoding and re-encoding keeps every value, dropping only the comments.

Given a sample.env file of:
```env
# comment
NAME="Mike Wazowski"
PATH=/usr/local/bin
GREETING="hello world"
QUOTED='say "hi"'

```
then
```bash
yq -p=dotenv -o=dotenv '.' sample.env
```
will output
```env
NAME='Mike Wazowski'
PATH=/usr/local/bin
GREETING='hello world'
QUOTED='say "hi"'
```

## Decode a line without a value
Given a sample.env file of:
```env
NAME
```
then an error is expected:
```
bad file 'sample.yml': line 1: expected KEY=value, got "NAME"
```

## Decode an unterminated quote
Given a sample.env file of:
```env
NAME="Mike
```
then an error is expected:
```
bad file 'sample.yml': line 1: unterminated double-quoted value "\"Mike"
```

## Encode a value dotenv cannot express
Given a sample.yml file of:
```yaml
mixed: >-
  it's "quoted"

```
then an error is expected:
```
cannot encode mixed: value "it's \"quoted\"" has both quote characters, which dotenv cannot express
```

