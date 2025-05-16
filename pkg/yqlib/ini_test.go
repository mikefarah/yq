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
		indent:       0,
		expected:     expectedSimpleINIOutput,
		scenarioType: "encode",
	},
	{
		description:  "Roundtrip INI: simple",
		input:        simpleINIInput,
		expected:     expectedSimpleINIOutput,
		scenarioType: "roundtrip",
		indent:       0,
	},
	{
		description:   "bad ini",
		input:         `[section\nkey = value`,
		expectedError: `bad file 'sample.yml': failed to parse INI content: unclosed section: [section\nkey = value`,
		scenarioType:  "decode-error",
	},
}

func documentRoundtripINIScenario(w *bufio.Writer, s formatScenario, indent int) {
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
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -p=ini -o=ini -I=%v '%v' sample.ini\n```\n", indent, expression))
	} else {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -p=ini -o=ini -I=%v sample.ini\n```\n", indent))
	}

	writeOrPanic(w, "will output\n")
	prefs := ConfiguredINIPreferences.Copy()
	prefs.Indent = indent

	// Pass prefs.Indent instead of prefs
	writeOrPanic(w, fmt.Sprintf("```ini\n%v```\n\n", mustProcessFormatScenario(s, NewINIDecoder(), NewINIEncoder(prefs.Indent))))
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
	prefs := ConfiguredINIPreferences.Copy()
	prefs.Indent = s.indent
	switch s.scenarioType {
	case "encode":
		// Pass prefs.Indent instead of prefs
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewINIEncoder(prefs.Indent)), s.description)
	case "decode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewINIDecoder(), NewYamlEncoder(ConfiguredYamlPreferences)), s.description)
	case "roundtrip":
		// Pass prefs.Indent instead of prefs
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewINIDecoder(), NewINIEncoder(prefs.Indent)), s.description)
	case "decode-error":
		// Pass prefs.Indent instead of prefs
		result, err := processFormatScenario(s, NewINIDecoder(), NewINIEncoder(prefs.Indent))
		if err == nil {
			t.Errorf("Expected error '%v' but it worked: %v", s.expectedError, result)
		} else {
			test.AssertResultComplexWithContext(t, s.expectedError, err.Error(), s.description)
		}
	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentINIScenario(t *testing.T, w *bufio.Writer, i interface{}) {
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
		documentRoundtripINIScenario(w, s, s.indent)
	case "decode-error":
		// Add handling for decode-error scenario type to prevent panic
		writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))
		if s.subdescription != "" {
			writeOrPanic(w, s.subdescription)
			writeOrPanic(w, "\n\n")
		}
		writeOrPanic(w, "Given a sample.ini file of:\n")
		writeOrPanic(w, fmt.Sprintf("```ini\n%v\n```\n", s.input))
		writeOrPanic(w, "then an error is expected:\n")
		writeOrPanic(w, fmt.Sprintf("```\n%v\n```\n\n", s.expectedError))
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

	if s.indent == 2 {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -o=ini '%v' sample.yml\n```\n", expression))
	} else {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -o=ini -I=%v '%v' sample.yml\n```\n", s.indent, expression))
	}
	writeOrPanic(w, "will output\n")
	prefs := ConfiguredINIPreferences.Copy()
	prefs.Indent = s.indent

	// Pass prefs.Indent instead of prefs
	writeOrPanic(w, fmt.Sprintf("```ini\n%v```\n\n", mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewINIEncoder(prefs.Indent))))
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
