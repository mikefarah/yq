package yqlib

import (
	"bufio"
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var inputXMLWithComments = `
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
var inputXMLWithCommentsWithSubChild = `
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

var expectedDecodeYamlWithSubChild = `# before cat
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
                +sweet: cool
            # in d after
        # in y after
    # in_cat_after
# after cat
`

var inputXMLWithCommentsWithArray = `
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

var expectedDecodeYamlWithArray = `# before cat
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
                +sweet: cool
              # in d after
            - # in d2 before
              z:
                +sweet: cool2
              # in d2 after
        # in y after
    # in_cat_after
# after cat
`

var expectedDecodeYamlWithComments = `# before cat
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

var expectedRoundtripXMLWithComments = `<!-- before cat --><cat><!-- in cat before -->
  <x>3<!-- multi
line comment 
for x --></x><!-- before y -->
  <y><!-- in y before
in d before -->
    <d>z<!-- in d after --></d><!-- in y after -->
  </y><!-- in_cat_after -->
</cat><!-- after cat -->
`

var yamlWithComments = `# above_cat
cat: # inline_cat
  # above_array
  array: # inline_array
    - val1 # inline_val1
    # above_val2
    - val2 # inline_val2
# below_cat
`

var expectedXMLWithComments = `<!-- above_cat inline_cat --><cat><!-- above_array inline_array -->
  <array>val1<!-- inline_val1 --></array>
  <array><!-- above_val2 -->val2<!-- inline_val2 --></array>
</cat><!-- below_cat -->
`

var inputXMLWithNamespacedAttr = `
<?xml version="1.0"?>
<map xmlns="some-namespace" xmlns:xsi="some-instance" xsi:schemaLocation="some-url">
</map>
`

var expectedYAMLWithNamespacedAttr = `map:
  +xmlns: some-namespace
  +xmlns:xsi: some-instance
  +some-instance:schemaLocation: some-url
`

var expectedYAMLWithRawNamespacedAttr = `map:
  +xmlns: some-namespace
  +xmlns:xsi: some-instance
  +xsi:schemaLocation: some-url
`

var xmlWithCustomDtd = `
<?xml version="1.0"?>
<!DOCTYPE root [
<!ENTITY writer "Blah.">
<!ENTITY copyright "Blah">
]>
<root>
    <item>&writer;&copyright;</item>
</root>`

var expectedDtd = `root:
    item: '&writer;&copyright;'
`

var xmlScenarios = []formatScenario{
	{
		description:    "Parse xml: simple",
		subdescription: "Notice how all the values are strings, see the next example on how you can fix that.",
		input:          "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<cat>\n  <says>meow</says>\n  <legs>4</legs>\n  <cute>true</cute>\n</cat>",
		expected:       "cat:\n    says: meow\n    legs: \"4\"\n    cute: \"true\"\n",
	},
	{
		description:    "Parse xml: number",
		subdescription: "All values are assumed to be strings when parsing XML, but you can use the `from_yaml` operator on all the strings values to autoparse into the correct type.",
		input:          "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<cat>\n  <says>meow</says>\n  <legs>4</legs>\n  <cute>true</cute>\n</cat>",
		expression:     " (.. | select(tag == \"!!str\")) |= from_yaml",
		expected:       "cat:\n    says: meow\n    legs: 4\n    cute: true\n",
	},
	{
		description:    "Parse xml: array",
		subdescription: "Consecutive nodes with identical xml names are assumed to be arrays.",
		input:          "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<animal>cat</animal>\n<animal>goat</animal>",
		expected:       "animal:\n    - cat\n    - goat\n",
	},
	{
		description:    "Parse xml: attributes",
		subdescription: "Attributes are converted to fields, with the default attribute prefix '+'. Use '--xml-attribute-prefix` to set your own.",
		input:          "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<cat legs=\"4\">\n  <legs>7</legs>\n</cat>",
		expected:       "cat:\n    +legs: \"4\"\n    legs: \"7\"\n",
	},
	{
		description:    "Parse xml: attributes with content",
		subdescription: "Content is added as a field, using the default content name of `+content`. Use `--xml-content-name` to set your own.",
		input:          "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<cat legs=\"4\">meow</cat>",
		expected:       "cat:\n    +content: meow\n    +legs: \"4\"\n",
	},
	{
		description:    "Parse xml: custom dtd",
		subdescription: "DTD entities are ignored.",
		input:          xmlWithCustomDtd,
		expected:       expectedDtd,
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
		description:  "Parse xml: keep attribute namespace",
		skipDoc:      false,
		input:        inputXMLWithNamespacedAttr,
		expected:     expectedYAMLWithNamespacedAttr,
		scenarioType: "decode-keep-ns",
	},
	{
		description:  "Parse xml: keep raw attribute namespace",
		skipDoc:      false,
		input:        inputXMLWithNamespacedAttr,
		expected:     expectedYAMLWithRawNamespacedAttr,
		scenarioType: "decode-raw-token",
	},
	{
		description:  "Encode xml: simple",
		input:        "cat: purrs",
		expected:     "<cat>purrs</cat>\n",
		scenarioType: "encode",
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
		input:          "cat:\n  +name: tiger\n  meows: true\n",
		expected:       "<cat name=\"tiger\">\n  <meows>true</meows>\n</cat>\n",
		scenarioType:   "encode",
	},
	{
		skipDoc:      true,
		input:        "cat:\n  ++name: tiger\n  meows: true\n",
		expected:     "<cat +name=\"tiger\">\n  <meows>true</meows>\n</cat>\n",
		scenarioType: "encode",
	},
	{
		description:    "Encode xml: attributes with content",
		subdescription: "Fields with the matching xml-content-name is assumed to be content.",
		input:          "cat:\n  +name: tiger\n  +content: cool\n",
		expected:       "<cat name=\"tiger\">cool</cat>\n",
		scenarioType:   "encode",
	},
	{
		description:    "Encode xml: comments",
		subdescription: "A best attempt is made to copy comments to xml.",
		input:          yamlWithComments,
		expected:       expectedXMLWithComments,
		scenarioType:   "encode",
	},
	{
		description:    "Round trip: with comments",
		subdescription: "A best effort is made, but comment positions and white space are not preserved perfectly.",
		input:          inputXMLWithComments,
		expected:       expectedRoundtripXMLWithComments,
		scenarioType:   "roundtrip",
	},
}

func testXMLScenario(t *testing.T, s formatScenario) {
	switch s.scenarioType {
	case "encode":
		test.AssertResultWithContext(t, s.expected, processFormatScenario(s, NewYamlDecoder(), NewXMLEncoder(2, "+", "+content")), s.description)
	case "roundtrip":
		test.AssertResultWithContext(t, s.expected, processFormatScenario(s, NewXMLDecoder("+", "+content", false, false, false), NewXMLEncoder(2, "+", "+content")), s.description)
	case "decode-keep-ns":
		test.AssertResultWithContext(t, s.expected, processFormatScenario(s, NewXMLDecoder("+", "+content", false, true, false), NewYamlEncoder(2, false, true, true)), s.description)
	case "decode-raw-token":
		test.AssertResultWithContext(t, s.expected, processFormatScenario(s, NewXMLDecoder("+", "+content", false, true, true), NewYamlEncoder(2, false, true, true)), s.description)
	default:
		test.AssertResultWithContext(t, s.expected, processFormatScenario(s, NewXMLDecoder("+", "+content", false, false, false), NewYamlEncoder(4, false, true, true)), s.description)
	}
}

func documentXMLScenario(t *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(formatScenario)

	if s.skipDoc {
		return
	}
	switch s.scenarioType {
	case "encode":
		documentXMLEncodeScenario(w, s)
	case "roundtrip":
		documentXMLRoundTripScenario(w, s)
	case "decode-keep-ns":
		documentXMLDecodeKeepNsScenario(w, s)
	case "decode-raw-token":
		documentXMLDecodeKeepNsRawTokenScenario(w, s)
	default:
		documentXMLDecodeScenario(w, s)
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
	if expression == "" {
		expression = "."
	}
	writeOrPanic(w, fmt.Sprintf("```bash\nyq -p=xml '%v' sample.xml\n```\n", expression))
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", processFormatScenario(s, NewXMLDecoder("+", "+content", false, false, false), NewYamlEncoder(2, false, true, true))))
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
	writeOrPanic(w, "```bash\nyq -p=xml -o=xml --xml-keep-namespace '.' sample.xml\n```\n")
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```xml\n%v```\n\n", processFormatScenario(s, NewXMLDecoder("+", "+content", false, true, false), NewXMLEncoder(2, "+", "+content"))))
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
	writeOrPanic(w, "```bash\nyq -p=xml -o=xml --xml-keep-namespace --xml-raw-token '.' sample.xml\n```\n")
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```xml\n%v```\n\n", processFormatScenario(s, NewXMLDecoder("+", "+content", false, true, true), NewXMLEncoder(2, "+", "+content"))))
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
	writeOrPanic(w, "```bash\nyq -o=xml '.' sample.yml\n```\n")
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```xml\n%v```\n\n", processFormatScenario(s, NewYamlDecoder(), NewXMLEncoder(2, "+", "+content"))))
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
	writeOrPanic(w, "```bash\nyq -p=xml -o=xml '.' sample.xml\n```\n")
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```xml\n%v```\n\n", processFormatScenario(s, NewXMLDecoder("+", "+content", false, false, false), NewXMLEncoder(2, "+", "+content"))))
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
