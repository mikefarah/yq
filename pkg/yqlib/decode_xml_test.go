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

type xmlScenario struct {
	inputXml       string
	expected       string
	description    string
	subdescription string
	skipDoc        bool
}

var xmlScenarios = []xmlScenario{
	{
		description: "Parse xml: simple",
		inputXml:    "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<cat>meow</cat>",
		expected:    "D0, P[], (doc)::cat: meow\n",
	},
	{
		description:    "Parse xml: array",
		subdescription: "Consecutive nodes with identical xml names are assumed to be arrays.",
		inputXml:       "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<animal>1</animal>\n<animal>2</animal>",
		expected:       "D0, P[], (doc)::animal:\n    - \"1\"\n    - \"2\"\n",
	},
	{
		description:    "Parse xml: attributes",
		subdescription: "Attributes are converted to fields, with the attribute prefix.",
		inputXml:       "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<cat legs=\"4\">\n  <legs>7</legs>\n</cat>",
		expected:       "D0, P[], (doc)::cat:\n    +legs: \"4\"\n    legs: \"7\"\n",
	},
	{
		description:    "Parse xml: attributes with content",
		subdescription: "Content is added as a field, using the content name",
		inputXml:       "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<cat legs=\"4\">meow</cat>",
		expected:       "D0, P[], (doc)::cat:\n    +content: meow\n    +legs: \"4\"\n",
	},
}

func testXmlScenario(t *testing.T, s *xmlScenario) {
	var actual = resultToString(t, decodeXml(t, s.inputXml))
	test.AssertResult(t, s.expected, actual)
}

func documentXmlScenario(t *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(xmlScenario)

	if s.skipDoc {
		return
	}
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.xml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```xml\n%v\n```\n", s.inputXml))

	writeOrPanic(w, "then\n")
	writeOrPanic(w, "```bash\nyq e sample.xml\n```\n")
	writeOrPanic(w, "will output\n")

	var output bytes.Buffer
	printer := NewPrinterWithSingleWriter(bufio.NewWriter(&output), YamlOutputFormat, true, false, 2, true)

	node := decodeXml(t, s.inputXml)

	printer.PrintResults(node.AsList())

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", output.String()))

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
