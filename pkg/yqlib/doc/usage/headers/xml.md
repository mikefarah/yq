# XML

At the moment, `yq` only supports decoding `xml` (into one of the other supported output formats).

As yaml does not have the concept of attributes, these are converted to regular fields with a prefix to prevent clobbering. Consecutive xml nodes with the same name are assumed to be arrays.

All values in XML are assumed to be strings - but you can use `from_yaml` to parse them into their correct types:


```
yq e -p=xml '.myNumberField |= from_yaml' my.xml
```
