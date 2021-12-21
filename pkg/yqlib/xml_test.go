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

func yamlToXml(sampleYaml string, indent int) string {
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	var encoder = NewXmlEncoder(writer, indent, "+")
	inputs, err := readDocuments(strings.NewReader(sampleYaml), "sample.yml", 0, NewYamlDecoder())
	if err != nil {
		panic(err)
	}
	node := inputs.Front().Value.(*CandidateNode).Node
	err = encoder.Encode(node)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	return strings.TrimSuffix(output.String(), "\n")
}

type xmlScenario struct {
	input          string
	expected       string
	description    string
	subdescription string
	skipDoc        bool
	encodeScenario bool
}

var xmlScenarios = []xmlScenario{
	{
		description: "Parse xml: simple",
		input:       "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<cat>meow</cat>",
		expected:    "D0, P[], (doc)::cat: meow\n",
	},
	{
		description:    "Parse xml: array",
		subdescription: "Consecutive nodes with identical xml names are assumed to be arrays.",
		input:          "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<animal>1</animal>\n<animal>2</animal>",
		expected:       "D0, P[], (doc)::animal:\n    - \"1\"\n    - \"2\"\n",
	},
	{
		description:    "Parse xml: attributes",
		subdescription: "Attributes are converted to fields, with the attribute prefix.",
		input:          "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<cat legs=\"4\">\n  <legs>7</legs>\n</cat>",
		expected:       "D0, P[], (doc)::cat:\n    +legs: \"4\"\n    legs: \"7\"\n",
	},
	{
		description:    "Parse xml: attributes with content",
		subdescription: "Content is added as a field, using the content name",
		input:          "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<cat legs=\"4\">meow</cat>",
		expected:       "D0, P[], (doc)::cat:\n    +content: meow\n    +legs: \"4\"\n",
	},
	{
		description:    "Encode xml: simple",
		input:          "cat: purrs",
		expected:       "<cat>purrs</cat>",
		encodeScenario: true,
	},
	{
		description:    "Encode xml: array",
		input:          "pets:\n  cat:\n    - purrs\n    - meows",
		expected:       "<pets>\n  <cat>purrs</cat>\n  <cat>meows</cat>\n</pets>",
		encodeScenario: true,
	},
	{
		description:    "Encode xml: attributes",
		subdescription: "Fields with the matching xml-attribute-prefix are assumed to be attributes.",
		input:          "cat:\n  +name: tiger\n  meows: true\n",
		expected:       "<cat name=\"tiger\">\n  <meows>true</meows>\n</cat>",
		encodeScenario: true,
	},
	{
		skipDoc:        true,
		input:          "cat:\n  ++name: tiger\n  meows: true\n",
		expected:       "<cat +name=\"tiger\">\n  <meows>true</meows>\n</cat>",
		encodeScenario: true,
	},
}

//encode

func testXmlScenario(t *testing.T, s *xmlScenario) {
	if s.encodeScenario {
		test.AssertResultWithContext(t, s.expected, yamlToXml(s.input, 2), s.description)
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
	if s.encodeScenario {
		documentXmlEncodeScenario(t, w, s)
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
	printer := NewPrinterWithSingleWriter(bufio.NewWriter(&output), YamlOutputFormat, true, false, 2, true)

	node := decodeXml(t, s.input)

	err := printer.PrintResults(node.AsList())
	if err != nil {
		t.Error(err)
		return
	}

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", output.String()))
}

func documentXmlEncodeScenario(t *testing.T, w *bufio.Writer, s xmlScenario) {
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

	writeOrPanic(w, fmt.Sprintf("```xml\n%v\n```\n\n", yamlToXml(s.input, 2)))
}

func TestXmlScenarios(t *testing.T) {
	for _, tt := range xmlScenarios {
		testXmlScenario(t, &tt)
	}
	genericScenarios := make([]interface{}, len(xmlScenarios))
	for i, s := range xmlScenarios {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "xml", genericScenarios, documentXmlScenario)
}
