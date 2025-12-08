//go:build !yq_noxml

package yqlib

import (
	"bufio"
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

const yamlInputWithProcInstAndHeadComment = `# cats
+p_xml: version="1.0"
this: is some xml`

const expectedXmlProcInstAndHeadComment = `<?xml version="1.0"?>
<!-- cats -->
<this>is some xml</this>
`

const xmlProcInstAndHeadCommentBlock = `<?xml version="1.0"?>
<!--
cats
-->
<this>is some xml</this>
`

const expectedYamlProcInstAndHeadCommentBlock = `#
# cats
#
+p_xml: version="1.0"
this: is some xml
`

const inputXMLWithComments = `
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
`
const inputXMLWithCommentsWithSubChild = `
<!-- before cat -->
<cat>
	<!-- in cat before -->
	<x>3<!-- multi
line comment 
for x --></x>
	<!-- before y -->
	<y>
		<!-- in y before -->
		<d><!-- in d before --><z sweet="cool"/><!-- in d after --></d>
		
		<!-- in y after -->
	</y>
	<!-- in_cat_after -->
</cat>
<!-- after cat -->
`

const expectedDecodeYamlWithSubChild = `# before cat
cat:
    # in cat before
    x: "3" # multi
    # line comment 
    # for x
    # before y

    y:
        # in y before
        d:
            # in d before
            z:
                +@sweet: cool
            # in d after
        # in y after
    # in_cat_after
# after cat
`

const inputXMLWithCommentsWithArray = `
<!-- before cat -->
<cat>
	<!-- in cat before -->
	<x>3<!-- multi
line comment 
for x --></x>
	<!-- before y -->
	<y>
		<!-- in y before -->
		<d><!-- in d before --><z sweet="cool"/><!-- in d after --></d>
        <d><!-- in d2 before --><z sweet="cool2"/><!-- in d2 after --></d>
		
		<!-- in y after -->
	</y>
	<!-- in_cat_after -->
</cat>
<!-- after cat -->
`

const expectedDecodeYamlWithArray = `# before cat
cat:
    # in cat before
    x: "3" # multi
    # line comment 
    # for x
    # before y

    y:
        # in y before
        d:
            - # in d before
              z:
                +@sweet: cool
              # in d after
            - # in d2 before
              z:
                +@sweet: cool2
              # in d2 after
        # in y after
    # in_cat_after
# after cat
`

const expectedDecodeYamlWithComments = `# before cat
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
`

const expectedRoundtripXMLWithComments = `<!-- before cat -->
<cat><!-- in cat before -->
  <x>3<!-- multi
line comment 
for x --></x><!-- before y -->
  <y><!-- in y before
in d before -->
    <d>z<!-- in d after --></d><!-- in y after -->
  </y><!-- in_cat_after -->
</cat><!-- after cat -->
`

const yamlWithComments = `#
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
`

const expectedXMLWithComments = `<!--
header comment
above_cat
-->
<!-- inline_cat -->
<cat><!-- above_array inline_array -->
  <array>val1<!-- inline_val1 --></array>
  <array><!-- above_val2 -->val2<!-- inline_val2 --></array>
</cat><!-- below_cat -->
`

const inputXMLWithNamespacedAttr = `<?xml version="1.0"?>
<map xmlns="some-namespace" xmlns:xsi="some-instance" xsi:schemaLocation="some-url">
  <item foo="bar">baz</item>
  <xsi:item>foobar</xsi:item>
</map>
`

const expectedYAMLWithNamespacedAttr = `+p_xml: version="1.0"
map:
  +@xmlns: some-namespace
  +@xmlns:xsi: some-instance
  +@xsi:schemaLocation: some-url
  item:
    +content: baz
    +@foo: bar
  xsi:item: foobar
`

const expectedYAMLWithRawNamespacedAttr = `+p_xml: version="1.0"
map:
  +@xmlns: some-namespace
  +@xmlns:xsi: some-instance
  +@xsi:schemaLocation: some-url
  item:
    +content: baz
    +@foo: bar
  xsi:item: foobar
`

const expectedYAMLWithoutRawNamespacedAttr = `+p_xml: version="1.0"
some-namespace:map:
  +@xmlns: some-namespace
  +@xmlns:xsi: some-instance
  +@some-instance:schemaLocation: some-url
  some-namespace:item:
    +content: baz
    +@foo: bar
  some-instance:item: foobar
`

const xmlWithCustomDtd = `
<?xml version="1.0"?>
<!DOCTYPE root [
<!ENTITY writer "Blah.">
<!ENTITY copyright "Blah">
]>
<root>
    <item>&writer;&copyright;</item>
</root>`

const expectedDtd = `<?xml version="1.0"?>
<!DOCTYPE root [
<!ENTITY writer "Blah.">
<!ENTITY copyright "Blah">
]>
<root>
  <item>&amp;writer;&amp;copyright;</item>
</root>
`

const expectedSkippedDtd = `<?xml version="1.0"?>
<root>
  <item>&amp;writer;&amp;copyright;</item>
</root>
`

const xmlWithProcInstAndDirectives = `<?xml version="1.0"?>
<!DOCTYPE config SYSTEM "/etc/iwatch/iwatch.dtd" >
<apple>
  <?coolioo version="1.0"?>
  <!CATYPE meow purr puss >
  <b>things</b>
</apple>
`

const yamlWithProcInstAndDirectives = `+p_xml: version="1.0"
+directive: 'DOCTYPE config SYSTEM "/etc/iwatch/iwatch.dtd" '
apple:
  +p_coolioo: version="1.0"
  +directive: 'CATYPE meow purr puss '
  b: things
`

const expectedXmlWithProcInstAndDirectives = `<?xml version="1.0"?>
<!DOCTYPE config SYSTEM "/etc/iwatch/iwatch.dtd" >
<apple><?coolioo version="1.0"?><!CATYPE meow purr puss >
  <b>things</b>
</apple>
`

var xmlScenarios = []formatScenario{
	{
		skipDoc:     true,
		description: "bad xml",
		input:       `<?xml version="1.0" encoding="UTF-8"?></Child></Root>`,
		expected:    "+p_xml: version=\"1.0\" encoding=\"UTF-8\"\n",
	},
	{
		skipDoc:  true,
		input:    "  <root>value<!-- comment--> </root>",
		expected: "root: value # comment\n",
	},
	{
		skipDoc:       true,
		input:         "value<root>value</root>",
		expectedError: "bad file 'sample.yml': invalid XML: Encountered chardata [value] outside of XML node",
		scenarioType:  "decode-error",
	},
	{
		skipDoc:  true,
		input:    "<root><!-- comment-->value</root>",
		expected: "# comment\nroot: value\n",
	},
	{
		skipDoc:  true,
		input:    "<root> <!-- comment--></root>",
		expected: "root: # comment\n",
	},
	{
		skipDoc:  true,
		input:    "<root>value<!-- comment-->anotherValue </root>",
		expected: "root:\n    # comment\n    - value\n    - anotherValue\n",
	},
	{
		skipDoc:  true,
		input:    "<root><cats><cat>quick</cat><cat>soft</cat><!-- kitty_comment--><cat>squishy</cat></cats></root>",
		expected: "root:\n    cats:\n        cat:\n            - quick\n            - soft\n            # kitty_comment\n\n            - squishy\n",
	},
	{
		description:  "ProcInst with head comment",
		skipDoc:      true,
		input:        yamlInputWithProcInstAndHeadComment,
		expected:     expectedXmlProcInstAndHeadComment,
		scenarioType: "encode",
	},
	{
		description:  "Scalar roundtrip",
		skipDoc:      true,
		input:        "<mike>cat</mike>",
		expression:   ".mike",
		expected:     "cat",
		scenarioType: "roundtrip",
	},
	{
		description:  "ProcInst with head comment round trip",
		skipDoc:      true,
		input:        expectedXmlProcInstAndHeadComment,
		expected:     expectedXmlProcInstAndHeadComment,
		scenarioType: "roundtrip",
	},
	{
		description:  "ProcInst with block head comment to yaml",
		skipDoc:      true,
		input:        xmlProcInstAndHeadCommentBlock,
		expected:     expectedYamlProcInstAndHeadCommentBlock,
		scenarioType: "decode",
	},
	{
		description:  "ProcInst with block head comment from yaml",
		skipDoc:      true,
		input:        expectedYamlProcInstAndHeadCommentBlock,
		expected:     xmlProcInstAndHeadCommentBlock,
		scenarioType: "encode",
	},
	{
		description:  "ProcInst with head comment round trip block",
		skipDoc:      true,
		input:        xmlProcInstAndHeadCommentBlock,
		expected:     xmlProcInstAndHeadCommentBlock,
		scenarioType: "roundtrip",
	},
	{
		description:    "Parse xml: simple",
		subdescription: "Notice how all the values are strings, see the next example on how you can fix that.",
		input:          "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<cat>\n  <says>meow</says>\n  <legs>4</legs>\n  <cute>true</cute>\n</cat>",
		expected:       "+p_xml: version=\"1.0\" encoding=\"UTF-8\"\ncat:\n    says: meow\n    legs: \"4\"\n    cute: \"true\"\n",
	},
	{
		description:    "Parse xml: number",
		subdescription: "All values are assumed to be strings when parsing XML, but you can use the `from_yaml` operator on all the strings values to autoparse into the correct type.",
		input:          "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<cat>\n  <says>meow</says>\n  <legs>4</legs>\n  <cute>true</cute>\n</cat>",
		expression:     " (.. | select(tag == \"!!str\")) |= from_yaml",
		expected:       "+p_xml: version=\"1.0\" encoding=\"UTF-8\"\ncat:\n    says: meow\n    legs: 4\n    cute: true\n",
	},
	{
		description:    "Parse xml: array",
		subdescription: "Consecutive nodes with identical xml names are assumed to be arrays.",
		input:          "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<animal>cat</animal>\n<animal>goat</animal>",
		expected:       "+p_xml: version=\"1.0\" encoding=\"UTF-8\"\nanimal:\n    - cat\n    - goat\n",
	},
	{
		description:    "Parse xml: force as an array",
		subdescription: "In XML, if your array has a single item, then yq doesn't know its an array. This is how you can consistently force it to be an array. This handles the 3 scenarios of having nothing in the array, having a single item and having multiple.",
		input:          "<zoo><animal>cat</animal></zoo>",
		expression:     ".zoo.animal |= ([] + .)",
		expected:       "zoo:\n    animal:\n        - cat\n",
	},
	{
		description: "Parse xml: force all as an array",
		input:       "<zoo><thing><frog>boing</frog></thing></zoo>",
		expression:  ".. |= [] + .",
		expected:    "- zoo:\n    - thing:\n        - frog:\n            - boing\n",
	},
	{
		description:    "Parse xml: attributes",
		subdescription: "Attributes are converted to fields, with the default attribute prefix '+'. Use '--xml-attribute-prefix` to set your own.",
		input:          "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<cat legs=\"4\">\n  <legs>7</legs>\n</cat>",
		expected:       "+p_xml: version=\"1.0\" encoding=\"UTF-8\"\ncat:\n    +@legs: \"4\"\n    legs: \"7\"\n",
	},
	{
		description:    "Parse xml: attributes with content",
		subdescription: "Content is added as a field, using the default content name of `+content`. Use `--xml-content-name` to set your own.",
		input:          "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<cat legs=\"4\">meow</cat>",
		expected:       "+p_xml: version=\"1.0\" encoding=\"UTF-8\"\ncat:\n    +content: meow\n    +@legs: \"4\"\n",
	},
	{
		description:    "Parse xml: content split between comments/children",
		subdescription: "Multiple content texts are collected into a sequence.",
		input:          "<root>  value  <!-- comment-->anotherValue <a>frog</a> cool!</root>",
		expected:       "root:\n    +content: # comment\n        - value\n        - anotherValue\n        - cool!\n    a: frog\n",
	},
	{
		description:    "Parse xml: custom dtd",
		subdescription: "DTD entities are processed as directives.",
		input:          xmlWithCustomDtd,
		expected:       expectedDtd,
		scenarioType:   "roundtrip",
	},
	{
		description:  "Roundtrip with name spaced attributes",
		skipDoc:      true,
		input:        inputXMLWithNamespacedAttr,
		expected:     inputXMLWithNamespacedAttr,
		scenarioType: "roundtrip",
	},
	{
		description:    "Parse xml: skip custom dtd",
		subdescription: "DTDs are directives, skip over directives to skip DTDs.",
		input:          xmlWithCustomDtd,
		expected:       expectedSkippedDtd,
		scenarioType:   "roundtrip-skip-directives",
	},
	{
		description:    "Parse xml: with comments",
		subdescription: "A best attempt is made to preserve comments.",
		input:          inputXMLWithComments,
		expected:       expectedDecodeYamlWithComments,
		scenarioType:   "decode",
	},
	{
		description:  "Empty doc",
		skipDoc:      true,
		input:        "",
		expected:     "\n",
		scenarioType: "decode",
	},
	{
		description:  "Empty single node",
		skipDoc:      true,
		input:        "<a/>",
		expected:     "a:\n",
		scenarioType: "decode",
	},
	{
		description:  "Empty close node",
		skipDoc:      true,
		input:        "<a></a>",
		expected:     "a:\n",
		scenarioType: "decode",
	},
	{
		description:  "Nested empty",
		skipDoc:      true,
		input:        "<a><b/></a>",
		expected:     "a:\n    b:\n",
		scenarioType: "decode",
	},
	{
		description:  "Parse xml: with comments subchild",
		skipDoc:      true,
		input:        inputXMLWithCommentsWithSubChild,
		expected:     expectedDecodeYamlWithSubChild,
		scenarioType: "decode",
	},
	{
		description:  "Parse xml: with comments array",
		skipDoc:      true,
		input:        inputXMLWithCommentsWithArray,
		expected:     expectedDecodeYamlWithArray,
		scenarioType: "decode",
	},
	{
		description:    "Parse xml: keep attribute namespace",
		subdescription: fmt.Sprintf(`Defaults to %v`, ConfiguredXMLPreferences.KeepNamespace),
		skipDoc:        false,
		input:          inputXMLWithNamespacedAttr,
		expected:       expectedYAMLWithNamespacedAttr,
		scenarioType:   "decode-keep-ns",
	},
	{
		description:  "Parse xml: keep raw attribute namespace",
		skipDoc:      true,
		input:        inputXMLWithNamespacedAttr,
		expected:     expectedYAMLWithRawNamespacedAttr,
		scenarioType: "decode-raw-token",
	},
	{
		description:    "Parse xml: keep raw attribute namespace",
		subdescription: fmt.Sprintf(`Defaults to %v`, ConfiguredXMLPreferences.UseRawToken),
		skipDoc:        false,
		input:          inputXMLWithNamespacedAttr,
		expected:       expectedYAMLWithoutRawNamespacedAttr,
		scenarioType:   "decode-raw-token-off",
	},
	{
		description:  "Encode xml: simple",
		input:        "cat: purrs",
		expected:     "<cat>purrs</cat>\n",
		scenarioType: "encode",
	},
	{
		description:  "includes map tags",
		skipDoc:      true,
		input:        "<cat>purrs</cat>\n",
		expression:   `tag`,
		expected:     "!!map\n",
		scenarioType: "decode",
	},
	{
		description:  "includes array tags",
		skipDoc:      true,
		input:        "<cat>purrs</cat><cat>purrs</cat>\n",
		expression:   `.cat | tag`,
		expected:     "!!seq\n",
		scenarioType: "decode",
	},
	{
		description:  "Encode xml: array",
		input:        "pets:\n  cat:\n    - purrs\n    - meows",
		expected:     "<pets>\n  <cat>purrs</cat>\n  <cat>meows</cat>\n</pets>\n",
		scenarioType: "encode",
	},
	{
		description:    "Encode xml: attributes",
		subdescription: "Fields with the matching xml-attribute-prefix are assumed to be attributes.",
		input:          "cat:\n  +@name: tiger\n  meows: true\n",
		expected:       "<cat name=\"tiger\">\n  <meows>true</meows>\n</cat>\n",
		scenarioType:   "encode",
	},
	{
		description:  "double prefix",
		skipDoc:      true,
		input:        "cat:\n  +@+@name: tiger\n  meows: true\n",
		expected:     "<cat +@name=\"tiger\">\n  <meows>true</meows>\n</cat>\n",
		scenarioType: "encode",
	},
	{
		description:   "arrays cannot be encoded",
		skipDoc:       true,
		input:         "[cat, dog, fish]",
		expectedError: "cannot encode !!seq to XML - only maps can be encoded",
		scenarioType:  "encode-error",
	},
	{
		description:   "arrays cannot be encoded - 2",
		skipDoc:       true,
		input:         "[cat, dog]",
		expectedError: "cannot encode !!seq to XML - only maps can be encoded",
		scenarioType:  "encode-error",
	},
	{
		description:    "Encode xml: attributes with content",
		subdescription: "Fields with the matching xml-content-name is assumed to be content.",
		input:          "cat:\n  +@name: tiger\n  +content: cool\n",
		expected:       "<cat name=\"tiger\">cool</cat>\n",
		scenarioType:   "encode",
	},
	{
		description:  "round trip multiline 1",
		skipDoc:      true,
		input:        "<x><!-- cats --></x>\n",
		expected:     "<x><!-- cats --></x>\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "round trip multiline 2",
		skipDoc:      true,
		input:        "<x><!--\n cats\n --></x>\n",
		expected:     "<x><!--\ncats\n--></x>\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "round trip multiline 3",
		skipDoc:      true,
		input:        "<x><!--\n\tcats\n --></x>\n",
		expected:     "<x><!--\n\tcats\n--></x>\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "round trip multiline 4",
		skipDoc:      true,
		input:        "<x><!--\n\tcats\n\tdogs\n--></x>\n",
		expected:     "<x><!--\n\tcats\n\tdogs\n--></x>\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "round trip multiline 5",
		skipDoc:      true, // pity spaces aren't kept atm.
		input:        "<x><!--\ncats\ndogs\n--></x>\n",
		expected:     "<x><!--\ncats\ndogs\n--></x>\n",
		scenarioType: "roundtrip",
	},
	{
		description:    "Encode xml: comments",
		subdescription: "A best attempt is made to copy comments to xml.",
		input:          yamlWithComments,
		expected:       expectedXMLWithComments,
		scenarioType:   "encode",
	},
	{
		description:    "Encode: doctype and xml declaration",
		subdescription: "Use the special xml names to add/modify proc instructions and directives.",
		input:          yamlWithProcInstAndDirectives,
		expected:       expectedXmlWithProcInstAndDirectives,
		scenarioType:   "encode",
	},
	{
		description:    "Round trip: with comments",
		subdescription: "A best effort is made, but comment positions and white space are not preserved perfectly.",
		input:          inputXMLWithComments,
		expected:       expectedRoundtripXMLWithComments,
		scenarioType:   "roundtrip",
	},
	{
		description:    "Roundtrip: with doctype and declaration",
		subdescription: "yq parses XML proc instructions and directives into nodes.\nUnfortunately the underlying XML parser loses whitespace information.",
		input:          xmlWithProcInstAndDirectives,
		expected:       expectedXmlWithProcInstAndDirectives,
		scenarioType:   "roundtrip",
	},
}

func testXMLScenario(t *testing.T, s formatScenario) {
	switch s.scenarioType {
	case "", "decode":
		yamlPrefs := ConfiguredYamlPreferences.Copy()
		yamlPrefs.Indent = 4
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewXMLDecoder(ConfiguredXMLPreferences), NewYamlEncoder(yamlPrefs)), s.description)
	case "encode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewXMLEncoder(ConfiguredXMLPreferences)), s.description)
	case "roundtrip":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewXMLDecoder(ConfiguredXMLPreferences), NewXMLEncoder(ConfiguredXMLPreferences)), s.description)
	case "decode-keep-ns":
		prefs := NewDefaultXmlPreferences()
		prefs.KeepNamespace = true
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewXMLDecoder(prefs), NewYamlEncoder(ConfiguredYamlPreferences)), s.description)
	case "decode-raw-token":
		prefs := NewDefaultXmlPreferences()
		prefs.UseRawToken = true
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewXMLDecoder(prefs), NewYamlEncoder(ConfiguredYamlPreferences)), s.description)
	case "decode-raw-token-off":
		prefs := NewDefaultXmlPreferences()
		prefs.UseRawToken = false
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewXMLDecoder(prefs), NewYamlEncoder(ConfiguredYamlPreferences)), s.description)
	case "roundtrip-skip-directives":
		prefs := NewDefaultXmlPreferences()
		prefs.SkipDirectives = true
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewXMLDecoder(prefs), NewXMLEncoder(prefs)), s.description)
	case "decode-error":
		result, err := processFormatScenario(s, NewXMLDecoder(NewDefaultXmlPreferences()), NewYamlEncoder(ConfiguredYamlPreferences))
		if err == nil {
			t.Errorf("Expected error '%v' but it worked: %v", s.expectedError, result)
		} else {
			test.AssertResultComplexWithContext(t, s.expectedError, err.Error(), s.description)
		}
	case "encode-error":
		result, err := processFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewXMLEncoder(NewDefaultXmlPreferences()))
		if err == nil {
			t.Errorf("Expected error '%v' but it worked: %v", s.expectedError, result)
		} else {
			test.AssertResultComplexWithContext(t, s.expectedError, err.Error(), s.description)
		}

	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentXMLScenario(_ *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(formatScenario)

	if s.skipDoc {
		return
	}
	switch s.scenarioType {
	case "", "decode":
		documentXMLDecodeScenario(w, s)
	case "encode":
		documentXMLEncodeScenario(w, s)
	case "roundtrip":
		documentXMLRoundTripScenario(w, s)
	case "decode-keep-ns":
		documentXMLDecodeKeepNsScenario(w, s)
	case "decode-raw-token-off":
		documentXMLDecodeKeepNsRawTokenScenario(w, s)
	case "roundtrip-skip-directives":
		documentXMLSkipDirectivesScenario(w, s)

	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentXMLDecodeScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.xml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```xml\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	expression := s.expression
	if s.expression != "" {
		expression = fmt.Sprintf(" '%v'", s.expression)
	}
	writeOrPanic(w, fmt.Sprintf("```bash\nyq -oy%v sample.xml\n```\n", expression))
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", mustProcessFormatScenario(s, NewXMLDecoder(ConfiguredXMLPreferences), NewYamlEncoder(ConfiguredYamlPreferences))))
}

func documentXMLDecodeKeepNsScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.xml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```xml\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	writeOrPanic(w, "```bash\nyq --xml-keep-namespace=false sample.xml\n```\n")
	writeOrPanic(w, "will output\n")
	prefs := NewDefaultXmlPreferences()
	prefs.KeepNamespace = false
	writeOrPanic(w, fmt.Sprintf("```xml\n%v```\n\n", mustProcessFormatScenario(s, NewXMLDecoder(prefs), NewXMLEncoder(prefs))))

	prefsWithout := NewDefaultXmlPreferences()
	prefs.KeepNamespace = true
	writeOrPanic(w, "instead of\n")
	writeOrPanic(w, fmt.Sprintf("```xml\n%v```\n\n", mustProcessFormatScenario(s, NewXMLDecoder(prefsWithout), NewXMLEncoder(prefsWithout))))
}

func documentXMLDecodeKeepNsRawTokenScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.xml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```xml\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	writeOrPanic(w, "```bash\nyq --xml-raw-token=false sample.xml\n```\n")
	writeOrPanic(w, "will output\n")

	prefs := NewDefaultXmlPreferences()
	prefs.UseRawToken = false

	writeOrPanic(w, fmt.Sprintf("```xml\n%v```\n\n", mustProcessFormatScenario(s, NewXMLDecoder(prefs), NewXMLEncoder(prefs))))

	prefsWithout := NewDefaultXmlPreferences()
	prefsWithout.UseRawToken = true

	writeOrPanic(w, "instead of\n")
	writeOrPanic(w, fmt.Sprintf("```xml\n%v```\n\n", mustProcessFormatScenario(s, NewXMLDecoder(prefsWithout), NewXMLEncoder(prefsWithout))))
}

func documentXMLEncodeScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.yml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```yaml\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	writeOrPanic(w, "```bash\nyq -o=xml sample.yml\n```\n")
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```xml\n%v```\n\n", mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewXMLEncoder(ConfiguredXMLPreferences))))
}

func documentXMLRoundTripScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.xml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```xml\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	writeOrPanic(w, "```bash\nyq sample.xml\n```\n")
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```xml\n%v```\n\n", mustProcessFormatScenario(s, NewXMLDecoder(ConfiguredXMLPreferences), NewXMLEncoder(ConfiguredXMLPreferences))))
}

func documentXMLSkipDirectivesScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.xml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```xml\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	writeOrPanic(w, "```bash\nyq --xml-skip-directives sample.xml\n```\n")
	writeOrPanic(w, "will output\n")
	prefs := NewDefaultXmlPreferences()
	prefs.SkipDirectives = true

	writeOrPanic(w, fmt.Sprintf("```xml\n%v```\n\n", mustProcessFormatScenario(s, NewXMLDecoder(prefs), NewXMLEncoder(prefs))))
}

func TestXMLScenarios(t *testing.T) {
	for _, tt := range xmlScenarios {
		testXMLScenario(t, tt)
	}
	genericScenarios := make([]interface{}, len(xmlScenarios))
	for i, s := range xmlScenarios {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "xml", genericScenarios, documentXMLScenario)
}
