//go:build !yq_nojson5

package yqlib

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var json5Scenarios = []formatScenario{
	{
		description: "Parse json5: comments, trailing commas, single quotes",
		input: `{
  // comment
  unquoted: 'single quoted',
  trailing: [1, 2,],
}
`,
		expected: `# comment
unquoted: single quoted
trailing:
  - 1
  - 2
`,
		scenarioType: "decode",
	},
	{
		description:  "Parse json5: multiline block comments",
		skipDoc:      true,
		scenarioType: "decode",
		input: `{
  /*
    multiline
    block comment
  */
  first: 1,
  second/* inline block */: 2,
  third: /* before value */ 3,
  fourth: [1, /* between elements */ 2,],
}
`,
		expected: `# multiline
# block comment
first: 1
second: 2 # inline block
third: 3
# before value
fourth:
  - 1
  # between elements
  - 2
`,
	},
	{
		description:  "Roundtrip json5: preserve comment placement",
		skipDoc:      true,
		scenarioType: "roundtrip",
		indent:       2,
		input: `{
  a/*k*/:/*v*/1/*after*/,
  b: 2 // end
}`,
		expected: `{
  "a" /* k */: /* v */ 1 /* after */,
  "b": 2 /* end */
}
`,
	},
	{
		description:  "Roundtrip json5: hex, Infinity, NaN",
		skipDoc:      true,
		input:        `{hex: 0x10, posInf: Infinity, negInf: -Infinity, not: NaN,}`,
		expected:     "{\"hex\":16,\"posInf\":Infinity,\"negInf\":-Infinity,\"not\":NaN}\n",
		indent:       0,
		scenarioType: "roundtrip",
	},
	{
		description:    "bad json5",
		skipDoc:        true,
		input:          `{a: 1,]`,
		expectedError:  `bad file 'sample.yml': json5: expected object key at line 1, column 7`,
		scenarioType:   "decode-error",
		subdescription: "json5 supports more syntax than json, but it still needs to be well-formed.",
	},
	{
		description:    "Parse json5: block comments normalisation (stars, tabs, trailing spaces)",
		skipDoc:        true,
		scenarioType:   "decode",
		input:          "{\n  /*\n\t *  hello \t \n\t ** world\t\n  */\n  a: 1,\n}\n",
		expected:       "# hello \t \n# * world\na: 1\n",
		subdescription: "Block comments are normalised and emitted as YAML comments.",
	},
	{
		description:  "Parse json5: scalars true/false/null with leading whitespace",
		skipDoc:      true,
		scenarioType: "decode",
		input:        " \n\ttrue\n",
		expected:     "true\n",
	},
	{
		description:  "Parse json5: scalar false",
		skipDoc:      true,
		scenarioType: "decode",
		input:        "false",
		expected:     "false\n",
	},
	{
		description:  "Parse json5: scalar null with extra newlines",
		skipDoc:      true,
		scenarioType: "decode",
		input:        "\n\nnull\n\n",
		expected:     "null\n",
	},
	{
		description:  "Parse json5: object trailing comma",
		skipDoc:      true,
		scenarioType: "decode",
		input:        "{a: 1,}\n",
		expected:     "a: 1\n",
	},
	{
		description:  "Parse json5: array trailing comma",
		skipDoc:      true,
		scenarioType: "decode",
		input:        "[1, 2,]\n",
		expected:     "- 1\n- 2\n",
	},
	{
		description:  "Parse json5: unicode escapes in strings",
		skipDoc:      true,
		scenarioType: "decode",
		input:        "{s: \"\\u0041\"}\n",
		expected:     "s: A\n",
	},
	{
		description:  "Parse json5: UTF-16 surrogate pair unicode escape in string",
		skipDoc:      true,
		scenarioType: "decode",
		input:        "{s: \"\\uD83D\\uDE00\"}\n",
		expected:     "s: \"\\U0001F600\"\n",
	},
	{
		description:           "Parse json5: invalid unicode escape in string (lowercase g)",
		skipDoc:               true,
		scenarioType:          "decode-error",
		input:                 "{s: \"\\u00ag\"}\n",
		expectedErrorContains: "invalid hex escape",
	},
	{
		description:           "Parse json5: invalid float (overflow)",
		skipDoc:               true,
		scenarioType:          "decode-error",
		input:                 "{a: 1e309}\n",
		expectedErrorContains: "invalid float number",
	},
	{
		description:           "Parse json5: invalid float exponent",
		skipDoc:               true,
		scenarioType:          "decode-error",
		input:                 "{a: 1e+}\n",
		expectedErrorContains: "invalid number exponent",
	},
	{
		description:           "Parse json5: invalid number with multiple periods",
		skipDoc:               true,
		scenarioType:          "decode-error",
		input:                 "{ip: 127.0.0.1}\n",
		expectedErrorContains: "expected ',' or '}' after object entry",
	},
	{
		description:           "Parse json5: extra colon after key",
		skipDoc:               true,
		scenarioType:          "decode-error",
		input:                 "{a:: 1}\n",
		expectedErrorContains: "unexpected character ':'",
	},
	{
		description:           "Parse json5: trailing comma with extra closing brace",
		skipDoc:               true,
		scenarioType:          "decode-error",
		input:                 "{a: 1,}}\n",
		expectedErrorContains: "unexpected character '}'",
	},
	{
		description:           "Parse json5: unterminated block comment at EOF",
		skipDoc:               true,
		scenarioType:          "decode-error",
		input:                 "{/* unterminated",
		expectedErrorContains: "unterminated block comment",
	},
	{
		description:           "Parse json5: unterminated escape sequence at EOF",
		skipDoc:               true,
		scenarioType:          "decode-error",
		input:                 "\"abc\\",
		expectedErrorContains: "unterminated escape sequence",
	},
	{
		description:           "Parse json5: unterminated unicode escape",
		skipDoc:               true,
		scenarioType:          "decode-error",
		input:                 "{s: \"\\u12\"}\n",
		expectedErrorContains: "invalid hex escape",
	},
	{
		description:           "Parse json5: bad hex digits in string escape",
		skipDoc:               true,
		scenarioType:          "decode-error",
		input:                 "{s: \"\\x0G\"}\n",
		expectedErrorContains: "invalid hex escape",
	},
	{
		description:           "Parse json5: bad hex digits in string escape (lowercase g)",
		skipDoc:               true,
		scenarioType:          "decode-error",
		input:                 "{s: \"\\xag\"}\n",
		expectedErrorContains: "invalid hex escape",
	},
	{
		description:           "Parse json5: bad hex digits in hex number",
		skipDoc:               true,
		scenarioType:          "decode-error",
		input:                 "{a: 0xG}\n",
		expectedErrorContains: "invalid hex number",
	},
	{
		description:           "Parse json5: bad hex digits in hex number (lowercase g)",
		skipDoc:               true,
		scenarioType:          "decode-error",
		input:                 "{a: 0xg}\n",
		expectedErrorContains: "invalid hex number",
	},
	{
		description:           "Parse json5: invalid identifier start via unicode escape",
		skipDoc:               true,
		scenarioType:          "decode-error",
		input:                 "{\\u0031abc: 1}\n",
		expectedErrorContains: "invalid identifier start",
	},
	{
		description:           "Parse json5: invalid identifier part via unicode escape",
		skipDoc:               true,
		scenarioType:          "decode-error",
		input:                 "{a\\u002D: 1}\n",
		expectedErrorContains: "invalid identifier part",
	},
	{
		description:           "Parse json5: invalid unicode escape in identifier",
		skipDoc:               true,
		scenarioType:          "decode-error",
		input:                 "{a\\u00ag: 1}\n",
		expectedErrorContains: "invalid hex escape",
	},
}

func testJSON5Scenario(t *testing.T, s formatScenario) {
	prefs := ConfiguredJSONPreferences.Copy()
	prefs.Indent = s.indent
	prefs.UnwrapScalar = false

	switch s.scenarioType {
	case "encode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewJSON5Encoder(prefs)), s.description)
	case "decode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewJSON5Decoder(), NewYamlEncoder(ConfiguredYamlPreferences)), s.description)
	case "roundtrip":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewJSON5Decoder(), NewJSON5Encoder(prefs)), s.description)
	case "decode-error":
		result, err := processFormatScenario(s, NewJSON5Decoder(), NewJSON5Encoder(prefs))
		if err == nil {
			t.Errorf("Expected error '%v' but it worked: %v", s.expectedError, result)
		} else {
			if s.expectedErrorContains != "" {
				if !strings.Contains(err.Error(), s.expectedErrorContains) {
					t.Errorf("%s: expected error containing %q, got %q", s.description, s.expectedErrorContains, err.Error())
				}
			} else {
				test.AssertResultComplexWithContext(t, s.expectedError, err.Error(), s.description)
			}
		}
	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentJSON5Scenario(_ *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(formatScenario)
	if s.skipDoc {
		return
	}

	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.json5 file of:\n")
	writeOrPanic(w, fmt.Sprintf("```json5\n%v\n```\n", s.input))
	writeOrPanic(w, "then\n")
	writeOrPanic(w, "```bash\nyq -P -p=json5 '.' sample.json5\n```\n")
	writeOrPanic(w, "will output\n")
	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", mustProcessFormatScenario(s, NewJSON5Decoder(), NewYamlEncoder(ConfiguredYamlPreferences))))
}

func TestJSON5Scenarios(t *testing.T) {
	for _, tt := range json5Scenarios {
		testJSON5Scenario(t, tt)
	}

	genericScenarios := make([]interface{}, len(json5Scenarios))
	for i, s := range json5Scenarios {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "json5", genericScenarios, documentJSON5Scenario)
}
