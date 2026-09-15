//go:build !yq_nobase64

package yqlib

import (
	"bufio"
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

const base64UrlEncodedSimple = "YSBzcGVjaWFsIHN0cmluZw=="
const base64UrlDecodedSimpleExtraSpaces = "\n " + base64UrlEncodedSimple + "  \n"
const base64UrlDecodedSimple = "a special string"

const base64UrlEncodedUTF8 = "V29ya3Mgd2l0aCBVVEYtMTYg8J-Yig=="
const base64UrlDecodedUTF8 = "Works with UTF-16 😊"

// contains both '-' and '_', where standard base64 would use '+' and '/'
const base64UrlEncodedUrlSafeChars = "Pj4tPz8_"
const base64UrlDecodedUrlSafeChars = ">>-???"

const base64UrlEncodedYaml = "YTogYXBwbGUK"
const base64UrlDecodedYaml = "a: apple\n"

const base64UrlEncodedEmpty = ""
const base64UrlDecodedEmpty = ""

const base64UrlMissingPadding = "Y2F0cw"
const base64UrlDecodedMissingPadding = "cats"

const base64UrlEncodedCats = "Y2F0cw=="
const base64UrlDecodedCats = "cats"

var base64UrlScenarios = []formatScenario{
	{
		skipDoc:      true,
		description:  "empty decode",
		input:        base64UrlEncodedEmpty,
		expected:     base64UrlDecodedEmpty + "\n",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "simple decode",
		input:        base64UrlEncodedSimple,
		expected:     base64UrlDecodedSimple + "\n",
		scenarioType: "decode",
	},
	{
		description:    "Decode base64url: simple",
		subdescription: "Decoded data is assumed to be a string.",
		input:          base64UrlEncodedSimple,
		expected:       base64UrlDecodedSimple + "\n",
		scenarioType:   "decode",
	},
	{
		description:    "Decode base64url: UTF-8",
		subdescription: "Base64url decoding supports UTF-8 encoded strings.",
		input:          base64UrlEncodedUTF8,
		expected:       base64UrlDecodedUTF8 + "\n",
		scenarioType:   "decode",
	},
	{
		description:    "Decode base64url: URL-safe characters",
		subdescription: "`-` and `_` are used in place of `+` and `/` (RFC 4648 §5).",
		input:          base64UrlEncodedUrlSafeChars,
		expected:       base64UrlDecodedUrlSafeChars + "\n",
		scenarioType:   "decode",
	},
	{
		skipDoc:      true,
		description:  "decode missing padding",
		input:        base64UrlMissingPadding,
		expected:     base64UrlDecodedMissingPadding + "\n",
		scenarioType: "decode",
	},
	{
		description:    "Decode with extra spaces",
		subdescription: "Extra leading/trailing whitespace is stripped",
		input:          base64UrlDecodedSimpleExtraSpaces,
		expected:       base64UrlDecodedSimple + "\n",
		scenarioType:   "decode",
	},
	{
		skipDoc:      true,
		description:  "decode with padding",
		input:        base64UrlEncodedCats,
		expected:     base64UrlDecodedCats + "\n",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "decode yaml document",
		input:        base64UrlEncodedYaml,
		expected:     base64UrlDecodedYaml + "\n",
		scenarioType: "decode",
	},
	{
		description:  "Encode base64url: string",
		input:        "\"" + base64UrlDecodedSimple + "\"",
		expected:     base64UrlEncodedSimple,
		scenarioType: "encode",
	},
	{
		description:    "Encode base64url: string from document",
		subdescription: "Extract a string field and encode it to base64url.",
		input:          "coolData: \"" + base64UrlDecodedSimple + "\"",
		expression:     ".coolData",
		expected:       base64UrlEncodedSimple,
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
		input:        "\"" + base64UrlDecodedUTF8 + "\"",
		expected:     base64UrlEncodedUTF8,
		scenarioType: "encode",
	},
	{
		skipDoc:      true,
		description:  "encode cats",
		input:        "\"" + base64UrlDecodedCats + "\"",
		expected:     base64UrlEncodedCats,
		scenarioType: "encode",
	},
	{
		description:  "Roundtrip: simple",
		skipDoc:      true,
		input:        base64UrlEncodedSimple,
		expected:     base64UrlEncodedSimple,
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: UTF-8",
		skipDoc:      true,
		input:        base64UrlEncodedUTF8,
		expected:     base64UrlEncodedUTF8,
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: missing padding",
		skipDoc:      true,
		input:        base64UrlMissingPadding,
		expected:     base64UrlEncodedCats,
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: empty",
		skipDoc:      true,
		input:        base64UrlEncodedEmpty,
		expected:     base64UrlEncodedEmpty,
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

func testBase64UrlScenario(t *testing.T, s formatScenario) {
	switch s.scenarioType {
	case "", "decode":
		yamlPrefs := ConfiguredYamlPreferences.Copy()
		yamlPrefs.Indent = 4
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewBase64URLDecoder(), NewYamlEncoder(yamlPrefs)), s.description)
	case "encode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewBase64URLEncoder()), s.description)
	case "roundtrip":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewBase64URLDecoder(), NewBase64URLEncoder()), s.description)
	case "encode-error":
		result, err := processFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewBase64URLEncoder())
		if err == nil {
			t.Errorf("Expected error '%v' but it worked: %v", s.expectedError, result)
		} else {
			test.AssertResultComplexWithContext(t, s.expectedError, err.Error(), s.description)
		}

	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentBase64UrlScenario(_ *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(formatScenario)

	if s.skipDoc {
		return
	}
	switch s.scenarioType {
	case "", "decode":
		documentBase64UrlDecodeScenario(w, s)
	case "encode":
		documentBase64UrlEncodeScenario(w, s)

	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentBase64UrlDecodeScenario(w *bufio.Writer, s formatScenario) {
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
	writeOrPanic(w, fmt.Sprintf("```bash\nyq -p=base64url -oy '%v' sample.txt\n```\n", expression))
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", mustProcessFormatScenario(s, NewBase64URLDecoder(), NewYamlEncoder(ConfiguredYamlPreferences))))
}

func documentBase64UrlEncodeScenario(w *bufio.Writer, s formatScenario) {
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
	writeOrPanic(w, fmt.Sprintf("```bash\nyq -o=base64url '%v' sample.yml\n```\n", expression))
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```\n%v```\n\n", mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewBase64URLEncoder())))
}

func TestBase64UrlScenarios(t *testing.T) {
	for _, tt := range base64UrlScenarios {
		testBase64UrlScenario(t, tt)
	}
	genericScenarios := make([]interface{}, len(base64UrlScenarios))
	for i, s := range base64UrlScenarios {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "base64url", genericScenarios, documentBase64UrlScenario)
}
