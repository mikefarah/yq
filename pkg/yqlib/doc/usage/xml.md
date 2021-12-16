# XML

At the moment, `yq` only supports decoding `xml` (into one of the other supported output formats).

As yaml does not have the concept of attributes, these are converted to regular fields with a prefix to prevent clobbering. Consecutive xml nodes with the same name are assumed to be arrays.

All values in XML are assumed to be strings - but you can use `from_yaml` to parse them into their correct types:


```
yq e -p=xml '.myNumberField |= from_yaml' my.xml
```

## Parse xml: simple
Given a sample.xml file of:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<cat>meow</cat>
```
then
```bash
yq e sample.xml
```
will output
```yaml
cat: meow
```

## Parse xml: array
Consecutive nodes with identical xml names are assumed to be arrays.

Given a sample.xml file of:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<animal>1</animal>
<animal>2</animal>
```
then
```bash
yq e sample.xml
```
will output
```yaml
animal:
  - "1"
  - "2"
```

## Parse xml: attributes
Attributes are converted to fields, with the attribute prefix.

Given a sample.xml file of:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<cat legs="4">
  <legs>7</legs>
</cat>
```
then
```bash
yq e sample.xml
```
will output
```yaml
cat:
  +legs: "4"
  legs: "7"
```

## Parse xml: attributes with content
Content is added as a field, using the content name

Given a sample.xml file of:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<cat legs="4">meow</cat>
```
then
```bash
yq e sample.xml
```
will output
```yaml
cat:
  +content: meow
  +legs: "4"
```

