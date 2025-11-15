# Base64

Encode and decode to and from Base64.

Base64 assumes [RFC4648](https://rfc-editor.org/rfc/rfc4648.html) encoding. Encoding and decoding both assume that the content is a UTF-8 string and not binary content.


See below for examples


## Decode base64: simple
Decoded data is assumed to be a string.

Given a sample.txt file of:
```
YSBzcGVjaWFsIHN0cmluZw==
```
then
```bash
yq -p=base64 -oy '.' sample.txt
```
will output
```yaml
a special string
```

## Decode base64: UTF-8
Base64 decoding supports UTF-8 encoded strings.

Given a sample.txt file of:
```
V29ya3Mgd2l0aCBVVEYtMTYg8J+Yig==
```
then
```bash
yq -p=base64 -oy '.' sample.txt
```
will output
```yaml
Works with UTF-16 😊
```

## Decode with extra spaces
Extra leading/trailing whitespace is stripped

Given a sample.txt file of:
```

 YSBzcGVjaWFsIHN0cmluZw==  

```
then
```bash
yq -p=base64 -oy '.' sample.txt
```
will output
```yaml
a special string
```

## Encode base64: string
Given a sample.yml file of:
```yaml
"a special string"
```
then
```bash
yq -o=base64 '.' sample.yml
```
will output
```
YSBzcGVjaWFsIHN0cmluZw==```

## Encode base64: string from document
Extract a string field and encode it to base64.

Given a sample.yml file of:
```yaml
coolData: "a special string"
```
then
```bash
yq -o=base64 '.coolData' sample.yml
```
will output
```
YSBzcGVjaWFsIHN0cmluZw==```

