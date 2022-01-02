package yqlib

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
	yaml "gopkg.in/yaml.v3"
)

func decodeXml(t *testing.T, xml string) *CandidateNode {
	decoder := NewXmlDecoder("+", "+content")

	decoder.Init(strings.NewReader(xml))

	node := &yaml.Node{}
	err := decoder.Decode(node)
	if err != nil {
		t.Error(err, "fail to decode", xml)
	}
	return &CandidateNode{Node: node}
}

func processScenario(s xmlScenario) string {
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	var encoder = NewXmlEncoder(2, "+", "+content")

	var decoder = NewYamlDecoder()
	if s.scenarioType == "roundtrip" {
		decoder = NewXmlDecoder("+", "+content")
	}

	inputs, err := readDocuments(strings.NewReader(s.input), "sample.yml", 0, decoder)
	if err != nil {
		panic(err)
	}
	node := inputs.Front().Value.(*CandidateNode).Node
	err = encoder.Encode(writer, node)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	return output.String()
}

type xmlScenario struct {
	input          string
	expected       string
	description    string
	subdescription string
	skipDoc        bool
	scenarioType   string
}

var inputXmlWithComments = `
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
var inputXmlWithCommentsWithSubChild = `
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

var expectedDecodeYamlWithSubChild = `D0, P[], (doc)::# before cat
cat:
    # in cat before
    x: "3" # multi
    # line comment 
    # for x
    # before y

    y:
        # in y before
        # in d before
        d:
            z:
                sweet: cool
								# in d after
        				# in y after
								# after cat
`

var inputXmlWithCommentsWithArray = `
<!-- before cat -->
<cat>
	<!-- in cat before -->
	<x>3<!-- multi
line comment 
for x --></x>
	<!-- before y -->
	<y>
		<!-- in y before -->
		<d><!-- in d before --><z/><!-- in d after --></d>
		
		<!-- in y after -->
	</y>
	<!-- in_cat_after -->
</cat>
<!-- after cat -->
`

var expectedDecodeYamlWithComments = `D0, P[], (doc)::# before cat
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

# after cat
`

var expectedRoundtripXmlWithComments = `<!-- before cat --><cat><!-- in cat before -->
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

var expectedXmlWithComments = `<!-- above_cat inline_cat --><cat><!-- above_array inline_array -->
  <array>val1<!-- inline_val1 --></array>
  <array><!-- above_val2 -->val2<!-- inline_val2 --></array>
</cat><!-- below_cat -->
`

var xmlScenarios = []xmlScenario{
	// {
	// 	description: "Parse xml: simple",
	// 	input:       "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<cat>meow</cat>",
	// 	expected:    "D0, P[], (doc)::cat: meow\n",
	// },
	// {
	// 	description:    "Parse xml: array",
	// 	subdescription: "Consecutive nodes with identical xml names are assumed to be arrays.",
	// 	input:          "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<animal>1</animal>\n<animal>2</animal>",
	// 	expected:       "D0, P[], (doc)::animal:\n    - \"1\"\n    - \"2\"\n",
	// },
	// {
	// 	description:    "Parse xml: attributes",
	// 	subdescription: "Attributes are converted to fields, with the attribute prefix.",
	// 	input:          "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<cat legs=\"4\">\n  <legs>7</legs>\n</cat>",
	// 	expected:       "D0, P[], (doc)::cat:\n    +legs: \"4\"\n    legs: \"7\"\n",
	// },
	// {
	// 	description:    "Parse xml: attributes with content",
	// 	subdescription: "Content is added as a field, using the content name",
	// 	input:          "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<cat legs=\"4\">meow</cat>",
	// 	expected:       "D0, P[], (doc)::cat:\n    +content: meow\n    +legs: \"4\"\n",
	// },
	// {
	// 	description:    "Parse xml: with comments",
	// 	subdescription: "A best attempt is made to preserve comments.",
	// 	input:          inputXmlWithComments,
	// 	expected:       expectedDecodeYamlWithComments,
	// 	scenarioType:   "decode",
	// },
	{
		description:  "Parse xml: with comments subchild",
		skipDoc:      true,
		input:        inputXmlWithCommentsWithSubChild,
		expected:     expectedDecodeYamlWithSubChild,
		scenarioType: "decode",
	},
	// {
	// 	description:  "Encode xml: simple",
	// 	input:        "cat: purrs",
	// 	expected:     "<cat>purrs</cat>\n",
	// 	scenarioType: "encode",
	// },
	// {
	// 	description:  "Encode xml: array",
	// 	input:        "pets:\n  cat:\n    - purrs\n    - meows",
	// 	expected:     "<pets>\n  <cat>purrs</cat>\n  <cat>meows</cat>\n</pets>\n",
	// 	scenarioType: "encode",
	// },
	// {
	// 	description:    "Encode xml: attributes",
	// 	subdescription: "Fields with the matching xml-attribute-prefix are assumed to be attributes.",
	// 	input:          "cat:\n  +name: tiger\n  meows: true\n",
	// 	expected:       "<cat name=\"tiger\">\n  <meows>true</meows>\n</cat>\n",
	// 	scenarioType:   "encode",
	// },
	// {
	// 	skipDoc:      true,
	// 	input:        "cat:\n  ++name: tiger\n  meows: true\n",
	// 	expected:     "<cat +name=\"tiger\">\n  <meows>true</meows>\n</cat>\n",
	// 	scenarioType: "encode",
	// },
	// {
	// 	description:    "Encode xml: attributes with content",
	// 	subdescription: "Fields with the matching xml-content-name is assumed to be content.",
	// 	input:          "cat:\n  +name: tiger\n  +content: cool\n",
	// 	expected:       "<cat name=\"tiger\">cool</cat>\n",
	// 	scenarioType:   "encode",
	// },
	// {
	// 	description:    "Encode xml: comments",
	// 	subdescription: "A best attempt is made to copy comments to xml.",
	// 	input:          yamlWithComments,
	// 	expected:       expectedXmlWithComments,
	// 	scenarioType:   "encode",
	// },
	// {
	// 	description:    "Round trip: with comments",
	// 	subdescription: "A best effort is made, but comment positions and white space are not preserved perfectly.",
	// 	input:          inputXmlWithComments,
	// 	expected:       expectedRoundtripXmlWithComments,
	// 	scenarioType:   "roundtrip",
	// },
}

func testXmlScenario(t *testing.T, s xmlScenario) {
	if s.scenarioType == "encode" || s.scenarioType == "roundtrip" {
		test.AssertResultWithContext(t, s.expected, processScenario(s), s.description)
	} else {
		var actual = resultToString(t, decodeXml(t, s.input))
		test.AssertResultWithContext(t, s.expected, actual, s.description)
	}
}

func documentXmlScenario(t *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(xmlScenario)

	if s.skipDoc {
		return
	}
	if s.scenarioType == "encode" {
		documentXmlEncodeScenario(w, s)
	} else {
		documentXmlDecodeScenario(t, w, s)
	}

}

func documentXmlDecodeScenario(t *testing.T, w *bufio.Writer, s xmlScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.xml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```xml\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	writeOrPanic(w, "```bash\nyq e -p=xml '.' sample.xml\n```\n")
	writeOrPanic(w, "will output\n")

	var output bytes.Buffer
	printer := NewSimpleYamlPrinter(bufio.NewWriter(&output), YamlOutputFormat, true, false, 2, true)

	node := decodeXml(t, s.input)

	err := printer.PrintResults(node.AsList())
	if err != nil {
		t.Error(err)
		return
	}

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", output.String()))
}

func documentXmlEncodeScenario(w *bufio.Writer, s xmlScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.yml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```yaml\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	writeOrPanic(w, "```bash\nyq e -o=xml '.' sample.yml\n```\n")
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```xml\n%v```\n\n", processScenario(s)))
}

func TestXmlScenarios(t *testing.T) {
	for _, tt := range xmlScenarios {
		testXmlScenario(t, tt)
	}
	genericScenarios := make([]interface{}, len(xmlScenarios))
	for i, s := range xmlScenarios {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "xml", genericScenarios, documentXmlScenario)
}
