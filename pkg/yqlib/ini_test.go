//go:build !yq_noini

package yqlib

import (
	"bufio"
	"fmt"
	"strings"
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

const quotedINIInput = `[section]
color_theme = "Default"
theme_background = "False"
`

const expectedQuotedINIOutput = `[section]
color_theme      = "Default"
theme_background = "False"
`

const colonKeyINIInput = `[this:section]
this:line = should really work
`

// expectedColonKeyDefaultYaml demonstrates the default "=:" key/value delimiters
// splitting "this:line" into key "this" with value "line = should really work",
// since ":" is treated the same as "=" by default.
const expectedColonKeyDefaultYaml = `this:section:
  this: line = should really work
`

// expectedColonKeyEqualsOnlyYaml shows the fix: with --ini-key-value-delimiters
// set to "=" only, "this:line" is kept intact as the key.
const expectedColonKeyEqualsOnlyYaml = `this:section:
  this:line: should really work
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
	{
		description:    "Parse INI: key with colon",
		subdescription: fmt.Sprintf("By default, the key/value delimiters are %q, so ':' is treated the same as '=' and the key is split at the first delimiter found. See the next example for how to avoid this using `--ini-key-value-delimiters`.", ConfiguredINIPreferences.KeyValueDelimiters),
		input:          colonKeyINIInput,
		expected:       expectedColonKeyDefaultYaml,
		scenarioType:   "decode",
	},
}

// iniPreserveQuotesPrefs returns INIPreferences with PreserveSurroundedQuote enabled.
func iniPreserveQuotesPrefs() INIPreferences {
	prefs := NewDefaultINIPreferences()
	prefs.PreserveSurroundedQuote = true
	return prefs
}

var iniPreserveQuotesScenarios = []formatScenario{
	{
		description:  "Roundtrip INI: preserve quotes",
		input:        quotedINIInput,
		expected:     expectedQuotedINIOutput,
		scenarioType: "roundtrip",
	},
}

// iniEqualsOnlyKeyValueDelimiterPrefs returns INIPreferences with KeyValueDelimiters
// set to "=" only, so that ":" is not treated as a key/value separator.
func iniEqualsOnlyKeyValueDelimiterPrefs() INIPreferences {
	prefs := NewDefaultINIPreferences()
	prefs.KeyValueDelimiters = "="
	return prefs
}

var iniKeyValueDelimitersScenarios = []formatScenario{
	{
		description:    "Parse INI: key with colon, using --ini-key-value-delimiters",
		subdescription: "Keys containing a colon (e.g. Mercurial config files) are parsed correctly when ':' is removed from the key/value delimiters.",
		input:          colonKeyINIInput,
		expected:       expectedColonKeyEqualsOnlyYaml,
		scenarioType:   "decode-key-value-delimiters",
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
	writeOrPanic(w, fmt.Sprintf("```ini\n%v```\n\n", mustProcessFormatScenario(s, NewINIDecoder(NewDefaultINIPreferences()), NewINIEncoder())))
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
	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", mustProcessFormatScenario(s, NewINIDecoder(NewDefaultINIPreferences()), NewYamlEncoder(ConfiguredYamlPreferences))))
}

func testINIScenario(t *testing.T, s formatScenario) {
	switch s.scenarioType {
	case "encode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewINIEncoder()), s.description)
	case "decode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewINIDecoder(NewDefaultINIPreferences()), NewYamlEncoder(ConfiguredYamlPreferences)), s.description)
	case "roundtrip":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewINIDecoder(NewDefaultINIPreferences()), NewINIEncoder()), s.description)
	case "decode-error":
		result, err := processFormatScenario(s, NewINIDecoder(NewDefaultINIPreferences()), NewINIEncoder())
		if err == nil {
			t.Errorf("Expected error '%v' but it worked: %v", s.expectedError, result)
		} else {
			test.AssertResultComplexWithContext(t, s.expectedError, err.Error(), s.description)
		}
	case "decode-key-value-delimiters":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewINIDecoder(iniEqualsOnlyKeyValueDelimiterPrefs()), NewYamlEncoder(ConfiguredYamlPreferences)), s.description)
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
	case "decode-key-value-delimiters":
		documentDecodeKeyValueDelimitersINIScenario(w, s)
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

func documentDecodeKeyValueDelimitersINIScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.ini file of:\n")
	writeOrPanic(w, fmt.Sprintf("```ini\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	writeOrPanic(w, "```bash\nyq -p=ini --ini-key-value-delimiters='=' sample.ini\n```\n")
	writeOrPanic(w, "will output\n")
	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", mustProcessFormatScenario(s, NewINIDecoder(iniEqualsOnlyKeyValueDelimiterPrefs()), NewYamlEncoder(ConfiguredYamlPreferences))))

	writeOrPanic(w, "instead of\n")
	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", mustProcessFormatScenario(s, NewINIDecoder(NewDefaultINIPreferences()), NewYamlEncoder(ConfiguredYamlPreferences))))
}

func TestINIDecoderInitResetsFinished(t *testing.T) {
	decoder := NewINIDecoder(NewDefaultINIPreferences())
	firstDocuments, err := readDocuments(strings.NewReader("[first]\nkey = value\n"), "first.ini", 0, decoder)
	if err != nil {
		t.Fatal(err)
	}
	test.AssertResult(t, 1, firstDocuments.Len())

	secondDocuments, err := readDocuments(strings.NewReader("[second]\nkey = value\n"), "second.ini", 1, decoder)
	if err != nil {
		t.Fatal(err)
	}
	test.AssertResult(t, 1, secondDocuments.Len())
}

func TestINIScenarios(t *testing.T) {
	for _, tt := range iniScenarios {
		testINIScenario(t, tt)
	}
	for _, tt := range iniKeyValueDelimitersScenarios {
		testINIScenario(t, tt)
	}
	genericScenarios := make([]interface{}, 0, len(iniScenarios)+len(iniKeyValueDelimitersScenarios))
	for _, s := range iniScenarios {
		genericScenarios = append(genericScenarios, s)
	}
	for _, s := range iniKeyValueDelimitersScenarios {
		genericScenarios = append(genericScenarios, s)
	}
	documentScenarios(t, "usage", "ini", genericScenarios, documentINIScenario)
}

func testINIPreserveQuotesScenario(t *testing.T, s formatScenario) {
	prefs := iniPreserveQuotesPrefs()
	switch s.scenarioType {
	case "roundtrip":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewINIDecoder(prefs), NewINIEncoder()), s.description)
	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func TestINIPreserveQuotesScenarios(t *testing.T) {
	for _, tt := range iniPreserveQuotesScenarios {
		testINIPreserveQuotesScenario(t, tt)
	}
}
