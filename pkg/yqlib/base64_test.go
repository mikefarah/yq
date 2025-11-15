//go:build !yq_nobase64

package yqlib

import (
	"bufio"
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

const base64EncodedSimple = "YSBzcGVjaWFsIHN0cmluZw=="
const base64DecodedSimpleExtraSpaces = "\n " + base64EncodedSimple + "  \n"
const base64DecodedSimple = "a special string"

const base64EncodedUTF8 = "V29ya3Mgd2l0aCBVVEYtMTYg8J+Yig=="
const base64DecodedUTF8 = "Works with UTF-16 😊"

const base64EncodedYaml = "YTogYXBwbGUK"
const base64DecodedYaml = "a: apple\n"

const base64EncodedEmpty = ""
const base64DecodedEmpty = ""

const base64MissingPadding = "Y2F0cw"
const base64DecodedMissingPadding = "cats"

const base64EncodedCats = "Y2F0cw=="
const base64DecodedCats = "cats"

var base64Scenarios = []formatScenario{
	{
		skipDoc:      true,
		description:  "empty decode",
		input:        base64EncodedEmpty,
		expected:     base64DecodedEmpty + "\n",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "simple decode",
		input:        base64EncodedSimple,
		expected:     base64DecodedSimple + "\n",
		scenarioType: "decode",
	},
	{
		description:    "Decode base64: simple",
		subdescription: "Decoded data is assumed to be a string.",
		input:          base64EncodedSimple,
		expected:       base64DecodedSimple + "\n",
		scenarioType:   "decode",
	},
	{
		description:    "Decode base64: UTF-8",
		subdescription: "Base64 decoding supports UTF-8 encoded strings.",
		input:          base64EncodedUTF8,
		expected:       base64DecodedUTF8 + "\n",
		scenarioType:   "decode",
	},
	{
		skipDoc:      true,
		description:  "decode missing padding",
		input:        base64MissingPadding,
		expected:     base64DecodedMissingPadding + "\n",
		scenarioType: "decode",
	},
	{

		description:    "Decode with extra spaces",
		subdescription: "Extra leading/trailing whitespace is stripped",
		input:          base64DecodedSimpleExtraSpaces,
		expected:       base64DecodedSimple + "\n",
		scenarioType:   "decode",
	},
	{
		skipDoc:      true,
		description:  "decode with padding",
		input:        base64EncodedCats,
		expected:     base64DecodedCats + "\n",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "decode yaml document",
		input:        base64EncodedYaml,
		expected:     base64DecodedYaml + "\n",
		scenarioType: "decode",
	},
	{
		description:  "Encode base64: string",
		input:        "\"" + base64DecodedSimple + "\"",
		expected:     base64EncodedSimple,
		scenarioType: "encode",
	},
	{
		description:    "Encode base64: string from document",
		subdescription: "Extract a string field and encode it to base64.",
		input:          "coolData: \"" + base64DecodedSimple + "\"",
		expression:     ".coolData",
		expected:       base64EncodedSimple,
		scenarioType:   "encode",
	},
	{
		skipDoc:      true,
		description:  "encode empty string",
		input:        "\"\"",
		expected:     "",
		scenarioType: "encode",
	},
	{
		skipDoc:      true,
		description:  "encode UTF-8 string",
		input:        "\"" + base64DecodedUTF8 + "\"",
		expected:     base64EncodedUTF8,
		scenarioType: "encode",
	},
	{
		skipDoc:      true,
		description:  "encode cats",
		input:        "\"" + base64DecodedCats + "\"",
		expected:     base64EncodedCats,
		scenarioType: "encode",
	},
	{
		description:  "Roundtrip: simple",
		skipDoc:      true,
		input:        base64EncodedSimple,
		expected:     base64EncodedSimple,
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: UTF-8",
		skipDoc:      true,
		input:        base64EncodedUTF8,
		expected:     base64EncodedUTF8,
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: missing padding",
		skipDoc:      true,
		input:        base64MissingPadding,
		expected:     base64EncodedCats,
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: empty",
		skipDoc:      true,
		input:        base64EncodedEmpty,
		expected:     base64EncodedEmpty,
		scenarioType: "roundtrip",
	},
	{
		description:   "Encode error: non-string",
		skipDoc:       true,
		input:         "123",
		expectedError: "cannot encode !!int as base64, can only operate on strings",
		scenarioType:  "encode-error",
	},
	{
		description:   "Encode error: array",
		skipDoc:       true,
		input:         "[1, 2, 3]",
		expectedError: "cannot encode !!seq as base64, can only operate on strings",
		scenarioType:  "encode-error",
	},
	{
		description:   "Encode error: map",
		skipDoc:       true,
		input:         "{b: c}",
		expectedError: "cannot encode !!map as base64, can only operate on strings",
		scenarioType:  "encode-error",
	},
}

func testBase64Scenario(t *testing.T, s formatScenario) {
	switch s.scenarioType {
	case "", "decode":
		yamlPrefs := ConfiguredYamlPreferences.Copy()
		yamlPrefs.Indent = 4
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewBase64Decoder(), NewYamlEncoder(yamlPrefs)), s.description)
	case "encode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewBase64Encoder()), s.description)
	case "roundtrip":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewBase64Decoder(), NewBase64Encoder()), s.description)
	case "encode-error":
		result, err := processFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewBase64Encoder())
		if err == nil {
			t.Errorf("Expected error '%v' but it worked: %v", s.expectedError, result)
		} else {
			test.AssertResultComplexWithContext(t, s.expectedError, err.Error(), s.description)
		}

	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentBase64Scenario(_ *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(formatScenario)

	if s.skipDoc {
		return
	}
	switch s.scenarioType {
	case "", "decode":
		documentBase64DecodeScenario(w, s)
	case "encode":
		documentBase64EncodeScenario(w, s)

	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentBase64DecodeScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.txt file of:\n")
	writeOrPanic(w, fmt.Sprintf("```\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	expression := s.expression
	if expression == "" {
		expression = "."
	}
	writeOrPanic(w, fmt.Sprintf("```bash\nyq -p=base64 -oy '%v' sample.txt\n```\n", expression))
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", mustProcessFormatScenario(s, NewBase64Decoder(), NewYamlEncoder(ConfiguredYamlPreferences))))
}

func documentBase64EncodeScenario(w *bufio.Writer, s formatScenario) {
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
	writeOrPanic(w, fmt.Sprintf("```bash\nyq -o=base64 '%v' sample.yml\n```\n", expression))
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```\n%v```\n\n", mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewBase64Encoder())))
}

func TestBase64Scenarios(t *testing.T) {
	for _, tt := range base64Scenarios {
		testBase64Scenario(t, tt)
	}
	genericScenarios := make([]interface{}, len(base64Scenarios))
	for i, s := range base64Scenarios {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "base64", genericScenarios, documentBase64Scenario)
}
