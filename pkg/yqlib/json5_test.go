//go:build !yq_nojson5

package yqlib

import (
	"bufio"
	"fmt"
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
			test.AssertResultComplexWithContext(t, s.expectedError, err.Error(), s.description)
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
