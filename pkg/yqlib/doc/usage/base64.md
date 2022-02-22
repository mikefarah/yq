
{% hint style="warning" %}
Note that versions prior to 4.18 require the 'eval/e' command to be specified.&#x20;

`yq e <exp> <file>`
{% endhint %}

## Decode Base64
Decoded data is assumed to be a string.

Given a sample.yml file of:
```yml
V29ya3Mgd2l0aCBVVEYtMTYg8J+Yig==
```
then
```bash
yq -p=props sample.properties
```
will output
```yaml
V29ya3Mgd2l0aCBVVEYtMTYg8J+Yig: =
```

