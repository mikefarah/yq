# Base64url

Encode and decode to and from the URL- and filename-safe Base64 variant.

Base64url assumes [RFC4648 §5](https://rfc-editor.org/rfc/rfc4648.html#section-5) encoding: `-` and `_` are used in place of `+` and `/`. Encoding and decoding both assume that the content is a UTF-8 string and not binary content.


See below for examples


## Decode base64url: simple
Decoded data is assumed to be a string.

Given a sample.txt file of:
```
YSBzcGVjaWFsIHN0cmluZw==
```
then
```bash
yq -p=base64url -oy '.' sample.txt
```
will output
```yaml
a special string
```

## Decode base64url: UTF-8
Base64url decoding supports UTF-8 encoded strings.

Given a sample.txt file of:
```
V29ya3Mgd2l0aCBVVEYtMTYg8J-Yig==
```
then
```bash
yq -p=base64url -oy '.' sample.txt
```
will output
```yaml
Works with UTF-16 😊
```

## Decode base64url: URL-safe characters
`-` and `_` are used in place of `+` and `/` (RFC 4648 §5).

Given a sample.txt file of:
```
Pj4tPz8_
```
then
```bash
yq -p=base64url -oy '.' sample.txt
```
will output
```yaml
>>-???
```

## Decode with extra spaces
Extra leading/trailing whitespace is stripped

Given a sample.txt file of:
```

 YSBzcGVjaWFsIHN0cmluZw==  

```
then
```bash
yq -p=base64url -oy '.' sample.txt
```
will output
```yaml
a special string
```

## Encode base64url: string
Given a sample.yml file of:
```yaml
"a special string"
```
then
```bash
yq -o=base64url '.' sample.yml
```
will output
```
YSBzcGVjaWFsIHN0cmluZw==```

## Encode base64url: string from document
Extract a string field and encode it to base64url.

Given a sample.yml file of:
```yaml
coolData: "a special string"
```
then
```bash
yq -o=base64url '.coolData' sample.yml
```
will output
```
YSBzcGVjaWFsIHN0cmluZw==```

