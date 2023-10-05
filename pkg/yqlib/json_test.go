//go:build !yq_nojson

package yqlib

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

const complexExpectYaml = `a: Easy! as one two three
b:
  c: 2
  d:
    - 3
    - 4
`

const sampleNdJson = `{"this": "is a multidoc json file"}
{"each": ["line is a valid json document"]}
{"a number": 4}
`

const sampleNdJsonKey = `{"a": "first", "b": "next", "ab": "last"}`

const expectedJsonKeysInOrder = `a: first
b: next
ab: last
`

const expectedNdJsonYaml = `this: is a multidoc json file
---
each:
  - line is a valid json document
---
a number: 4
`

const expectedRoundTripSampleNdJson = `{"this":"is a multidoc json file"}
{"each":["line is a valid json document"]}
{"a number":4}
`

const expectedUpdatedMultilineJson = `{"this":"is a multidoc json file"}
{"each":["line is a valid json document","cool"]}
{"a number":4}
`

const sampleMultiLineJson = `{
	"this": "is a multidoc json file"
}
{
	"it": [
		"has",
		"consecutive",
		"json documents"
	]
}
{
	"a number": 4
}
`

const roundTripMultiLineJson = `{
  "this": "is a multidoc json file"
}
{
  "it": [
    "has",
    "consecutive",
    "json documents"
  ]
}
{
  "a number": 4
}
`

var jsonScenarios = []formatScenario{
	{
		description:  "array null",
		skipDoc:      true,
		input:        "[null]",
		scenarioType: "roundtrip-ndjson",
		expected:     "[null]\n",
	},
	{
		description:  "set tags",
		skipDoc:      true,
		input:        "[{}]",
		expression:   `[.. | type]`,
		scenarioType: "roundtrip-ndjson",
		expected:     "[\"!!seq\",\"!!map\"]\n",
	},
	{
		description:    "Parse json: simple",
		subdescription: "JSON is a subset of yaml, so all you need to do is prettify the output",
		input:          `{"cat": "meow"}`,
		scenarioType:   "decode-ndjson",
		expected:       "cat: meow\n",
	},
	{
		skipDoc:      true,
		description:  "Parse json: simple: key",
		input:        `{"cat": "meow"}`,
		expression:   ".cat | key",
		expected:     "\"cat\"\n",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "Parse json: simple: parent",
		input:        `{"cat": "meow"}`,
		expression:   ".cat | parent",
		expected:     "{\"cat\":\"meow\"}\n",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "Parse json: simple: path",
		input:        `{"cat": "meow"}`,
		expression:   ".cat | path",
		expected:     "[\"cat\"]\n",
		scenarioType: "decode",
	},
	{
		description:   "bad json",
		skipDoc:       true,
		input:         `{"a": 1 b": 2}`,
		expectedError: `bad file 'sample.yml': json: string of object unexpected end of JSON input`,
		scenarioType:  "decode-error",
	},
	{
		description:    "Parse json: complex",
		subdescription: "JSON is a subset of yaml, so all you need to do is prettify the output",
		input:          `{"a":"Easy! as one two three","b":{"c":2,"d":[3,4]}}`,
		expected:       complexExpectYaml,
		scenarioType:   "decode-ndjson",
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
	{
		description:    "Roundtrip NDJSON",
		subdescription: "Unfortunately the json encoder strips leading spaces of values.",
		input:          sampleNdJson,
		expected:       expectedRoundTripSampleNdJson,
		scenarioType:   "roundtrip-ndjson",
	},
	{
		description:    "Roundtrip multi-document JSON",
		subdescription: "The NDJSON parser can also handle multiple multi-line json documents in a single file!",
		input:          sampleMultiLineJson,
		expected:       roundTripMultiLineJson,
		scenarioType:   "roundtrip-multi",
	},
	{
		description:    "Update a specific document in a multi-document json",
		subdescription: "Documents are indexed by the `documentIndex` or `di` operator.",
		input:          sampleNdJson,
		expected:       expectedUpdatedMultilineJson,
		expression:     `(select(di == 1) | .each ) += "cool"`,
		scenarioType:   "roundtrip-ndjson",
	},
	{
		description:    "Find and update a specific document in a multi-document json",
		subdescription: "Use expressions as you normally would.",
		input:          sampleNdJson,
		expected:       expectedUpdatedMultilineJson,
		expression:     `(select(has("each")) | .each ) += "cool"`,
		scenarioType:   "roundtrip-ndjson",
	},
	{
		description:  "Decode NDJSON",
		input:        sampleNdJson,
		expected:     expectedNdJsonYaml,
		scenarioType: "decode-ndjson",
	},
	{
		description:  "Decode NDJSON, maintain key order",
		skipDoc:      true,
		input:        sampleNdJsonKey,
		expected:     expectedJsonKeysInOrder,
		scenarioType: "decode-ndjson",
	},
	{
		description:  "numbers",
		skipDoc:      true,
		input:        "[3, 3.0, 3.1, -1, 999999, 1000000, 1000001, 1.1]",
		expected:     "- 3\n- 3\n- 3.1\n- -1\n- 999999\n- 1000000\n- 1000001\n- 1.1\n",
		scenarioType: "decode-ndjson",
	},
	{
		description:  "number single",
		skipDoc:      true,
		input:        "3",
		expected:     "3\n",
		scenarioType: "decode-ndjson",
	},
	{
		description:  "empty string",
		skipDoc:      true,
		input:        `""`,
		expected:     "\n",
		scenarioType: "decode-ndjson",
	},
	{
		description:  "strings",
		skipDoc:      true,
		input:        `["", "cat"]`,
		expected:     "- \"\"\n- cat\n",
		scenarioType: "decode-ndjson",
	},
	{
		description:  "null",
		skipDoc:      true,
		input:        `null`,
		expected:     "null\n",
		scenarioType: "decode-ndjson",
	},
	{
		description:  "booleans",
		skipDoc:      true,
		input:        `[true, false]`,
		expected:     "- true\n- false\n",
		scenarioType: "decode-ndjson",
	},
}

func documentRoundtripNdJsonScenario(w *bufio.Writer, s formatScenario, indent int) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.json file of:\n")
	writeOrPanic(w, fmt.Sprintf("```json\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")

	expression := s.expression
	if expression != "" {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -p=json -o=json -I=%v '%v' sample.json\n```\n", indent, expression))
	} else {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -p=json -o=json -I=%v sample.json\n```\n", indent))
	}

	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", mustProcessFormatScenario(s, NewJSONDecoder(), NewJSONEncoder(indent, false, false))))
}

func documentDecodeNdJsonScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.json file of:\n")
	writeOrPanic(w, fmt.Sprintf("```json\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")

	expression := s.expression
	if expression != "" {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -p=json '%v' sample.json\n```\n", expression))
	} else {
		writeOrPanic(w, "```bash\nyq -p=json sample.json\n```\n")
	}

	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", mustProcessFormatScenario(s, NewJSONDecoder(), NewYamlEncoder(s.indent, false, ConfiguredYamlPreferences))))
}

func decodeJSON(t *testing.T, jsonString string) *CandidateNode {
	docs, err := readDocument(jsonString, "sample.json", 0)

	if err != nil {
		t.Error(err)
		return nil
	}

	exp, err := getExpressionParser().ParseExpression(PrettyPrintExp)

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

func testJSONScenario(t *testing.T, s formatScenario) {
	switch s.scenarioType {
	case "encode", "decode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewJSONEncoder(s.indent, false, false)), s.description)
	case "decode-ndjson":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewJSONDecoder(), NewYamlEncoder(2, false, ConfiguredYamlPreferences)), s.description)
	case "roundtrip-ndjson":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewJSONDecoder(), NewJSONEncoder(0, false, false)), s.description)
	case "roundtrip-multi":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewJSONDecoder(), NewJSONEncoder(2, false, false)), s.description)
	case "decode-error":
		result, err := processFormatScenario(s, NewJSONDecoder(), NewJSONEncoder(2, false, false))
		if err == nil {
			t.Errorf("Expected error '%v' but it worked: %v", s.expectedError, result)
		} else {
			test.AssertResultComplexWithContext(t, s.expectedError, err.Error(), s.description)
		}
	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentJSONDecodeScenario(t *testing.T, w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.json file of:\n")
	writeOrPanic(w, fmt.Sprintf("```json\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	writeOrPanic(w, "```bash\nyq -P '.' sample.json\n```\n")
	writeOrPanic(w, "will output\n")

	var output bytes.Buffer
	printer := NewSimpleYamlPrinter(bufio.NewWriter(&output), YamlOutputFormat, true, false, 2, true)

	node := decodeJSON(t, s.input)

	err := printer.PrintResults(node.AsList())
	if err != nil {
		t.Error(err)
		return
	}

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", output.String()))
}

func documentJSONScenario(t *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(formatScenario)
	if s.skipDoc {
		return
	}
	switch s.scenarioType {
	case "":
		documentJSONDecodeScenario(t, w, s)
	case "encode":
		documentJSONEncodeScenario(w, s)
	case "decode-ndjson":
		documentDecodeNdJsonScenario(w, s)
	case "roundtrip-ndjson":
		documentRoundtripNdJsonScenario(w, s, 0)
	case "roundtrip-multi":
		documentRoundtripNdJsonScenario(w, s, 2)

	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentJSONEncodeScenario(w *bufio.Writer, s formatScenario) {
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
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -o=json '%v' sample.yml\n```\n", expression))
	} else {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -o=json -I=%v '%v' sample.yml\n```\n", s.indent, expression))
	}
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```json\n%v```\n\n", mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewJSONEncoder(s.indent, false, false))))
}

func TestJSONScenarios(t *testing.T) {
	for _, tt := range jsonScenarios {
		testJSONScenario(t, tt)
	}
	genericScenarios := make([]interface{}, len(jsonScenarios))
	for i, s := range jsonScenarios {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "convert", genericScenarios, documentJSONScenario)
}
