//go:build !yq_noini

package yqlib

import (
	"bufio"
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

const simpleINIInput = `[section]
key = value
`

const expectedSimpleINIOutput = `[section]
key = value
`

const expectedSimpleINIYaml = `section:
  key: value
`

var iniScenarios = []formatScenario{
	{
		description:  "Parse INI: simple",
		input:        simpleINIInput,
		scenarioType: "decode",
		expected:     expectedSimpleINIYaml,
	},
	{
		description:  "Encode INI: simple",
		input:        `section: {key: value}`,
		expected:     expectedSimpleINIOutput,
		scenarioType: "encode",
	},
	{
		description:  "Roundtrip INI: simple",
		input:        simpleINIInput,
		expected:     expectedSimpleINIOutput,
		scenarioType: "roundtrip",
	},
	{
		description:   "bad ini",
		input:         `[section\nkey = value`,
		expectedError: `bad file 'sample.yml': failed to parse INI content: unclosed section: [section\nkey = value`,
		scenarioType:  "decode-error",
	},
}

func documentRoundtripINIScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.ini file of:\n")
	writeOrPanic(w, fmt.Sprintf("```ini\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")

	expression := s.expression
	if expression != "" {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -p=ini -o=ini '%v' sample.ini\n```\n", expression))
	} else {
		writeOrPanic(w, "```bash\nyq -p=ini -o=ini sample.ini\n```\n")
	}

	writeOrPanic(w, "will output\n")
	writeOrPanic(w, fmt.Sprintf("```ini\n%v```\n\n", mustProcessFormatScenario(s, NewINIDecoder(), NewINIEncoder())))
}

func documentDecodeINIScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.ini file of:\n")
	writeOrPanic(w, fmt.Sprintf("```ini\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")

	expression := s.expression
	if expression != "" {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -p=ini '%v' sample.ini\n```\n", expression))
	} else {
		writeOrPanic(w, "```bash\nyq -p=ini sample.ini\n```\n")
	}

	writeOrPanic(w, "will output\n")
	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", mustProcessFormatScenario(s, NewINIDecoder(), NewYamlEncoder(ConfiguredYamlPreferences))))
}

func testINIScenario(t *testing.T, s formatScenario) {
	switch s.scenarioType {
	case "encode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewINIEncoder()), s.description)
	case "decode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewINIDecoder(), NewYamlEncoder(ConfiguredYamlPreferences)), s.description)
	case "roundtrip":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewINIDecoder(), NewINIEncoder()), s.description)
	case "decode-error":
		result, err := processFormatScenario(s, NewINIDecoder(), NewINIEncoder())
		if err == nil {
			t.Errorf("Expected error '%v' but it worked: %v", s.expectedError, result)
		} else {
			test.AssertResultComplexWithContext(t, s.expectedError, err.Error(), s.description)
		}
	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentINIScenario(_ *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(formatScenario)
	if s.skipDoc {
		return
	}
	switch s.scenarioType {
	case "encode":
		documentINIEncodeScenario(w, s)
	case "decode":
		documentDecodeINIScenario(w, s)
	case "roundtrip":
		documentRoundtripINIScenario(w, s)
	case "decode-error":
		documentDecodeErrorINIScenario(w, s)
	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentINIEncodeScenario(w *bufio.Writer, s formatScenario) {
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

	writeOrPanic(w, fmt.Sprintf("```bash\nyq -o=ini '%v' sample.yml\n```\n", expression))

	writeOrPanic(w, "will output\n")
	writeOrPanic(w, fmt.Sprintf("```ini\n%v```\n\n", mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewINIEncoder())))
}

func documentDecodeErrorINIScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.ini file of:\n")
	writeOrPanic(w, fmt.Sprintf("```ini\n%v\n```\n", s.input))

	writeOrPanic(w, "then an error is expected:\n")
	writeOrPanic(w, fmt.Sprintf("```\n%v\n```\n\n", s.expectedError))
}

func TestINIScenarios(t *testing.T) {
	for _, tt := range iniScenarios {
		testINIScenario(t, tt)
	}
	genericScenarios := make([]interface{}, len(iniScenarios))
	for i, s := range iniScenarios {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "convert", genericScenarios, documentINIScenario)
}
