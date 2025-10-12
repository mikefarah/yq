# XML

Encode and decode to and from XML. Whitespace is not conserved for round trips - but the order of the fields are.

Consecutive xml nodes with the same name are assumed to be arrays.

XML content data, attributes processing instructions and directives are all created as plain fields. 

This can be controlled by:

| Flag | Default |Sample XML | 
| -- | -- |  -- |
 | `--xml-attribute-prefix` | `+` (changing to `+@` soon) | Legs in ```<cat legs="4"/>``` |  
 |  `--xml-content-name` | `+content` | Meow in ```<cat>Meow <fur>true</true></cat>``` |
 | `--xml-directive-name` | `+directive` | ```<!DOCTYPE config system "blah">``` |
 | `--xml-proc-inst-prefix` | `+p_` |  ```<?xml version="1"?>``` |


{% hint style="warning" %}
Default Attribute Prefix will be changing in v4.30!
In order to avoid name conflicts (e.g. having an attribute named "content" will create a field that clashes with the default content name of "+content") the attribute prefix will be changing to "+@".

This will affect users that have not set their own prefix and are not roundtripping XML changes.

{% endhint %}

## Encoder / Decoder flag options

In addition to the above flags, there are the following xml encoder/decoder options controlled by flags:

| Flag | Default | Description |
| -- | -- | -- |
| `--xml-strict-mode` | false | Strict mode enforces the requirements of the XML specification. When switched off the parser allows input containing common mistakes. See [the Golang xml decoder ](https://pkg.go.dev/encoding/xml#Decoder) for more details.| 
| `--xml-keep-namespace` | true | Keeps the namespace of attributes |
| `--xml-raw-token` | true |  Does not verify that start and end elements match and does not translate name space prefixes to their corresponding URLs. |
| `--xml-skip-proc-inst` | false | Skips over processing instructions, e.g. `<?xml version="1"?>` |
| `--xml-skip-directives` | false | Skips over directives, e.g. ```<!DOCTYPE config system "blah">``` |


See below for examples

## Parse xml: simple
Notice how all the values are strings, see the next example on how you can fix that.

Given a sample.xml file of:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<cat>
  <says>meow</says>
  <legs>4</legs>
  <cute>true</cute>
</cat>
```
then
```bash
yq -oy '.' sample.xml
```
will output
```yaml
+p_xml: version="1.0" encoding="UTF-8"
cat:
  says: meow
  legs: "4"
  cute: "true"
```

## Parse xml: number
All values are assumed to be strings when parsing XML, but you can use the `from_yaml` operator on all the strings values to autoparse into the correct type.

Given a sample.xml file of:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<cat>
  <says>meow</says>
  <legs>4</legs>
  <cute>true</cute>
</cat>
```
then
```bash
yq -oy ' (.. | select(tag == "!!str")) |= from_yaml' sample.xml
```
will output
```yaml
+p_xml: version="1.0" encoding="UTF-8"
cat:
  says: meow
  legs: 4
  cute: true
```

## Parse xml: array
Consecutive nodes with identical xml names are assumed to be arrays.

Given a sample.xml file of:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<animal>cat</animal>
<animal>goat</animal>
```
then
```bash
yq -oy '.' sample.xml
```
will output
```yaml
+p_xml: version="1.0" encoding="UTF-8"
animal:
  - cat
  - goat
```

## Parse xml: force as an array
In XML, if your array has a single item, then yq doesn't know its an array. This is how you can consistently force it to be an array. This handles the 3 scenarios of having nothing in the array, having a single item and having multiple.

Given a sample.xml file of:
```xml
<zoo><animal>cat</animal></zoo>
```
then
```bash
yq -oy '.zoo.animal |= ([] + .)' sample.xml
```
will output
```yaml
zoo:
  animal:
    - cat
```

## Parse xml: force all as an array
Given a sample.xml file of:
```xml
<zoo><thing><frog>boing</frog></thing></zoo>
```
then
```bash
yq -oy '.. |= [] + .' sample.xml
```
will output
```yaml
- zoo:
    - thing:
        - frog:
            - boing
```

## Parse xml: attributes
Attributes are converted to fields, with the default attribute prefix '+'. Use '--xml-attribute-prefix` to set your own.

Given a sample.xml file of:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<cat legs="4">
  <legs>7</legs>
</cat>
```
then
```bash
yq -oy '.' sample.xml
```
will output
```yaml
+p_xml: version="1.0" encoding="UTF-8"
cat:
  +@legs: "4"
  legs: "7"
```

## Parse xml: attributes with content
Content is added as a field, using the default content name of `+content`. Use `--xml-content-name` to set your own.

Given a sample.xml file of:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<cat legs="4">meow</cat>
```
then
```bash
yq -oy '.' sample.xml
```
will output
```yaml
+p_xml: version="1.0" encoding="UTF-8"
cat:
  +content: meow
  +@legs: "4"
```

## Parse xml: content split between comments/children
Multiple content texts are collected into a sequence.

Given a sample.xml file of:
```xml
<root>  value  <!-- comment-->anotherValue <a>frog</a> cool!</root>
```
then
```bash
yq -oy '.' sample.xml
```
will output
```yaml
root:
  +content: # comment
    - value
    - anotherValue
    - cool!
  a: frog
```

## Parse xml: custom dtd
DTD entities are processed as directives.

Given a sample.xml file of:
```xml

<?xml version="1.0"?>
<!DOCTYPE root [
<!ENTITY writer "Blah.">
<!ENTITY copyright "Blah">
]>
<root>
    <item>&writer;&copyright;</item>
</root>
```
then
```bash
yq '.' sample.xml
```
will output
```xml
<?xml version="1.0"?>
<!DOCTYPE root [
<!ENTITY writer "Blah.">
<!ENTITY copyright "Blah">
]>
<root>
  <item>&amp;writer;&amp;copyright;</item>
</root>
```

## Parse xml: skip custom dtd
DTDs are directives, skip over directives to skip DTDs.

Given a sample.xml file of:
```xml

<?xml version="1.0"?>
<!DOCTYPE root [
<!ENTITY writer "Blah.">
<!ENTITY copyright "Blah">
]>
<root>
    <item>&writer;&copyright;</item>
</root>
```
then
```bash
yq --xml-skip-directives '.' sample.xml
```
will output
```xml
<?xml version="1.0"?>
<root>
  <item>&amp;writer;&amp;copyright;</item>
</root>
```

## Parse xml: with comments
A best attempt is made to preserve comments.

Given a sample.xml file of:
```xml

<!-- before cat -->
<cat>
	<!-- in cat before -->
	<x>3<!-- multi
line comment 
for x --></x>
	<!-- before y -->
	<y>
		<!-- in y before -->
		<d><!-- in d before -->z<!-- in d after --></d>
		
		<!-- in y after -->
	</y>
	<!-- in_cat_after -->
</cat>
<!-- after cat -->

```
then
```bash
yq -oy '.' sample.xml
```
will output
```yaml
# before cat
cat:
  # in cat before
  x: "3" # multi
  # line comment 
  # for x
  # before y

  y:
    # in y before
    # in d before
    d: z # in d after
    # in y after
  # in_cat_after
# after cat
```

## Parse xml: keep attribute namespace
Defaults to true

Given a sample.xml file of:
```xml
<?xml version="1.0"?>
<map xmlns="some-namespace" xmlns:xsi="some-instance" xsi:schemaLocation="some-url">
  <item foo="bar">baz</item>
  <xsi:item>foobar</xsi:item>
</map>

```
then
```bash
yq --xml-keep-namespace=false '.' sample.xml
```
will output
```xml
<?xml version="1.0"?>
<map xmlns="some-namespace" xsi="some-instance" schemaLocation="some-url">
  <item foo="bar">baz</item>
  <item>foobar</item>
</map>
```

instead of
```xml
<?xml version="1.0"?>
<map xmlns="some-namespace" xmlns:xsi="some-instance" xsi:schemaLocation="some-url">
  <item foo="bar">baz</item>
  <xsi:item>foobar</xsi:item>
</map>
```

## Parse xml: keep raw attribute namespace
Defaults to true

Given a sample.xml file of:
```xml
<?xml version="1.0"?>
<map xmlns="some-namespace" xmlns:xsi="some-instance" xsi:schemaLocation="some-url">
  <item foo="bar">baz</item>
  <xsi:item>foobar</xsi:item>
</map>

```
then
```bash
yq --xml-raw-token=false '.' sample.xml
```
will output
```xml
<?xml version="1.0"?>
<some-namespace:map xmlns="some-namespace" xmlns:xsi="some-instance" some-instance:schemaLocation="some-url">
  <some-namespace:item foo="bar">baz</some-namespace:item>
  <some-instance:item>foobar</some-instance:item>
</some-namespace:map>
```

instead of
```xml
<?xml version="1.0"?>
<map xmlns="some-namespace" xmlns:xsi="some-instance" xsi:schemaLocation="some-url">
  <item foo="bar">baz</item>
  <xsi:item>foobar</xsi:item>
</map>
```

## Encode xml: simple
Given a sample.yml file of:
```yaml
cat: purrs
```
then
```bash
yq -o=xml sample.yml
```
will output
```xml
<cat>purrs</cat>
```

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
yq -o=xml sample.yml
```
will output
```xml
<pets>
  <cat>purrs</cat>
  <cat>meows</cat>
</pets>
```

## Encode xml: attributes
Fields with the matching xml-attribute-prefix are assumed to be attributes.

Given a sample.yml file of:
```yaml
cat:
  +@name: tiger
  meows: true

```
then
```bash
yq -o=xml sample.yml
```
will output
```xml
<cat name="tiger">
  <meows>true</meows>
</cat>
```

## Encode xml: attributes with content
Fields with the matching xml-content-name is assumed to be content.

Given a sample.yml file of:
```yaml
cat:
  +@name: tiger
  +content: cool

```
then
```bash
yq -o=xml sample.yml
```
will output
```xml
<cat name="tiger">cool</cat>
```

## Encode xml: comments
A best attempt is made to copy comments to xml.

Given a sample.yml file of:
```yaml
#
# header comment
# above_cat
#
cat: # inline_cat
  # above_array
  array: # inline_array
    - val1 # inline_val1
    # above_val2
    - val2 # inline_val2
# below_cat

```
then
```bash
yq -o=xml sample.yml
```
will output
```xml
<!--
header comment
above_cat
-->
<!-- inline_cat -->
<cat><!-- above_array inline_array -->
  <array>val1<!-- inline_val1 --></array>
  <array><!-- above_val2 -->val2<!-- inline_val2 --></array>
</cat><!-- below_cat -->
```

## Encode: doctype and xml declaration
Use the special xml names to add/modify proc instructions and directives.

Given a sample.yml file of:
```yaml
+p_xml: version="1.0"
+directive: 'DOCTYPE config SYSTEM "/etc/iwatch/iwatch.dtd" '
apple:
  +p_coolioo: version="1.0"
  +directive: 'CATYPE meow purr puss '
  b: things

```
then
```bash
yq -o=xml sample.yml
```
will output
```xml
<?xml version="1.0"?>
<!DOCTYPE config SYSTEM "/etc/iwatch/iwatch.dtd" >
<apple><?coolioo version="1.0"?><!CATYPE meow purr puss >
  <b>things</b>
</apple>
```

## Round trip: with comments
A best effort is made, but comment positions and white space are not preserved perfectly.

Given a sample.xml file of:
```xml

<!-- before cat -->
<cat>
	<!-- in cat before -->
	<x>3<!-- multi
line comment 
for x --></x>
	<!-- before y -->
	<y>
		<!-- in y before -->
		<d><!-- in d before -->z<!-- in d after --></d>
		
		<!-- in y after -->
	</y>
	<!-- in_cat_after -->
</cat>
<!-- after cat -->

```
then
```bash
yq '.' sample.xml
```
will output
```xml
<!-- before cat -->
<cat><!-- in cat before -->
  <x>3<!-- multi
line comment 
for x --></x><!-- before y -->
  <y><!-- in y before
in d before -->
    <d>z<!-- in d after --></d><!-- in y after -->
  </y><!-- in_cat_after -->
</cat><!-- after cat -->
```

## Roundtrip: with doctype and declaration
yq parses XML proc instructions and directives into nodes.
Unfortunately the underlying XML parser loses whitespace information.

Given a sample.xml file of:
```xml
<?xml version="1.0"?>
<!DOCTYPE config SYSTEM "/etc/iwatch/iwatch.dtd" >
<apple>
  <?coolioo version="1.0"?>
  <!CATYPE meow purr puss >
  <b>things</b>
</apple>

```
then
```bash
yq '.' sample.xml
```
will output
```xml
<?xml version="1.0"?>
<!DOCTYPE config SYSTEM "/etc/iwatch/iwatch.dtd" >
<apple><?coolioo version="1.0"?><!CATYPE meow purr puss >
  <b>things</b>
</apple>
```

