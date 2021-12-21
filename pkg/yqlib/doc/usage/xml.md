# XML

Encode and decode to and from XML. Whitespace is not conserved for round trips - but the order of the fields are.

As yaml does not have the concept of attributes, xml attributes are converted to regular fields with a prefix to prevent clobbering. This defaults to "+", use the `--xml-attribute-prefix` to change.

Consecutive xml nodes with the same name are assumed to be arrays.

All values in XML are assumed to be strings - but you can use `from_yaml` to parse them into their correct types:


```
yq e -p=xml '.myNumberField |= from_yaml' my.xml
```


XML nodes that have attributes then plain content, e.g:

```xml
<cat name="tiger">meow</cat>
```

The content of the node will be set as a field in the map with the key "+content". Use the `--xml-content-name` flag to change this.

## Parse xml: simple
Given a sample.xml file of:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<cat>meow</cat>
```
then
```bash
yq e -p=xml '.' sample.xml
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
yq e -p=xml '.' sample.xml
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
yq e -p=xml '.' sample.xml
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
yq e -p=xml '.' sample.xml
```
will output
```yaml
cat:
  +content: meow
  +legs: "4"
```

## Encode xml: simple
Given a sample.yml file of:
```yaml
cat: purrs
```
then
```bash
yq e -o=xml '.' sample.yml
```
will output
```xml
<cat>purrs</cat>```

## Encode xml: array
Given a sample.yml file of:
```yaml
pets:
  cat:
    - purrs
    - meows
```
then
```bash
yq e -o=xml '.' sample.yml
```
will output
```xml
<pets>
  <cat>purrs</cat>
  <cat>meows</cat>
</pets>```

## Encode xml: attributes
Fields with the matching xml-attribute-prefix are assumed to be attributes.

Given a sample.yml file of:
```yaml
cat:
  +name: tiger
  meows: true

```
then
```bash
yq e -o=xml '.' sample.yml
```
will output
```xml
<cat name="tiger">
  <meows>true</meows>
</cat>```

## Encode xml: attributes with content
Fields with the matching xml-content-name is assumed to be content.

Given a sample.yml file of:
```yaml
cat:
  +name: tiger
  +content: cool

```
then
```bash
yq e -o=xml '.' sample.yml
```
will output
```xml
<cat name="tiger">cool</cat>```

