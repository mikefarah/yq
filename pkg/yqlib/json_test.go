package yqlib

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var complexExpectYaml = `D0, P[], (!!map)::a: Easy! as one two three
b:
    c: 2
    d:
        - 3
        - 4
`

var jsonScenarios = []formatScenario{
	{
		description:    "Parse json: simple",
		subdescription: "JSON is a subset of yaml, so all you need to do is prettify the output",
		input:          `{"cat": "meow"}`,
		expected:       "D0, P[], (!!map)::cat: meow\n",
	},
	{
		description:    "Parse json: complex",
		subdescription: "JSON is a subset of yaml, so all you need to do is prettify the output",
		input:          `{"a":"Easy! as one two three","b":{"c":2,"d":[3,4]}}`,
		expected:       complexExpectYaml,
	},
	{
		description:  "Encode json: simple",
		input:        `cat: meow`,
		indent:       2,
		expected:     "{\n  \"cat\": \"meow\"\n}\n",
		scenarioType: "encode",
	},
	{
		description:  "Encode json: simple - in one line",
		input:        `cat: meow # this is a comment, and it will be dropped.`,
		indent:       0,
		expected:     "{\"cat\":\"meow\"}\n",
		scenarioType: "encode",
	},
	{
		description:  "Encode json: comments",
		input:        `cat: meow # this is a comment, and it will be dropped.`,
		indent:       2,
		expected:     "{\n  \"cat\": \"meow\"\n}\n",
		scenarioType: "encode",
	},
	{
		description:    "Encode json: anchors",
		subdescription: "Anchors are dereferenced",
		input:          "cat: &ref meow\nanotherCat: *ref",
		indent:         2,
		expected:       "{\n  \"cat\": \"meow\",\n  \"anotherCat\": \"meow\"\n}\n",
		scenarioType:   "encode",
	},
	{
		description:    "Encode json: multiple results",
		subdescription: "Each matching node is converted into a json doc. This is best used with 0 indent (json document per line)",
		input:          `things: [{stuff: cool}, {whatever: cat}]`,
		expression:     `.things[]`,
		indent:         0,
		expected:       "{\"stuff\":\"cool\"}\n{\"whatever\":\"cat\"}\n",
		scenarioType:   "encode",
	},
}

func decodeJson(t *testing.T, jsonString string) *CandidateNode {
	docs, err := readDocumentWithLeadingContent(jsonString, "sample.json", 0)

	if err != nil {
		t.Error(err)
		return nil
	}

	exp, err := NewExpressionParser().ParseExpression(PrettyPrintExp)

	if err != nil {
		t.Error(err)
		return nil
	}

	context, err := NewDataTreeNavigator().GetMatchingNodes(Context{MatchingNodes: docs}, exp)

	if err != nil {
		t.Error(err)
		return nil
	}

	return context.MatchingNodes.Front().Value.(*CandidateNode)
}

func testJsonScenario(t *testing.T, s formatScenario) {
	if s.scenarioType == "encode" || s.scenarioType == "roundtrip" {
		test.AssertResultWithContext(t, s.expected, processJsonScenario(s), s.description)
	} else {
		var actual = resultToString(t, decodeJson(t, s.input))
		test.AssertResultWithContext(t, s.expected, actual, s.description)
	}
}

func processJsonScenario(s formatScenario) string {

	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	var encoder = NewJsonEncoder(s.indent)

	var decoder = NewYamlDecoder()

	inputs, err := readDocuments(strings.NewReader(s.input), "sample.yml", 0, decoder)
	if err != nil {
		panic(err)
	}

	expression := s.expression
	if expression == "" {
		expression = "."
	}

	exp, err := NewExpressionParser().ParseExpression(expression)

	if err != nil {
		panic(err)
	}

	context, err := NewDataTreeNavigator().GetMatchingNodes(Context{MatchingNodes: inputs}, exp)

	if err != nil {
		panic(err)
	}

	printer := NewPrinter(encoder, NewSinglePrinterWriter(writer))
	err = printer.PrintResults(context.MatchingNodes)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	return output.String()

}

func documentJsonDecodeScenario(t *testing.T, w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.json file of:\n")
	writeOrPanic(w, fmt.Sprintf("```json\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	writeOrPanic(w, "```bash\nyq e -P '.' sample.json\n```\n")
	writeOrPanic(w, "will output\n")

	var output bytes.Buffer
	printer := NewSimpleYamlPrinter(bufio.NewWriter(&output), YamlOutputFormat, true, false, 2, true)

	node := decodeJson(t, s.input)

	err := printer.PrintResults(node.AsList())
	if err != nil {
		t.Error(err)
		return
	}

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", output.String()))
}

func documentJsonScenario(t *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(formatScenario)

	if s.skipDoc {
		return
	}
	if s.scenarioType == "encode" {
		documentJsonEncodeScenario(w, s)
	} else {
		documentJsonDecodeScenario(t, w, s)
	}
}

func documentJsonEncodeScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.yml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```yaml\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")

	expression := s.expression
	if expression == "" {
		expression = "."
	}

	if s.indent == 2 {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq e -o=json '%v' sample.yml\n```\n", expression))
	} else {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq e -o=json -I=%v '%v' sample.yml\n```\n", s.indent, expression))
	}
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```json\n%v```\n\n", processJsonScenario(s)))
}

func TestJsonScenarios(t *testing.T) {
	for _, tt := range jsonScenarios {
		testJsonScenario(t, tt)
	}
	genericScenarios := make([]interface{}, len(jsonScenarios))
	for i, s := range jsonScenarios {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "convert", genericScenarios, documentJsonScenario)
}
