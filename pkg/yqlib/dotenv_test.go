//go:build !yq_nodotenv

package yqlib

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var dotenvScenarios = []formatScenario{
	{
		description:    "Decode a dotenv file",
		subdescription: "Comments and blank lines are dropped and an `export` prefix is accepted. The value is a string, so `1` stays quoted in YAML output.",
		input: "" +
			"# a comment" + "\n" +
			"NAME=Mike Wazowski" + "\n" +
			"export ROLE=scarer" + "\n" +
			"\n" +
			"EYES=1 # inline comment" + "\n",
		expected: "" +
			"NAME: Mike Wazowski" + "\n" +
			"ROLE: scarer" + "\n" +
			"EYES: \"1\"" + "\n",
		scenarioType: "decode",
	},
	{
		description:    "Decode quoted values",
		subdescription: "`\\n` is the only escape sequence, so `\\t` stays two characters. A single-quoted value is taken literally and a quote ends a double-quoted value.",
		input: "" +
			`NEWLINE="line one\nline two"` + "\n" +
			`TAB="a\tb"` + "\n" +
			`SINGLE='line one\nline two'` + "\n" +
			`HASH="a # b"` + "\n" +
			`TRUNCATED="say \"hi\""` + "\n",
		expected: "" +
			"NEWLINE: |-" + "\n" +
			"  line one" + "\n" +
			"  line two" + "\n" +
			"TAB: a\\tb" + "\n" +
			"SINGLE: line one\\nline two" + "\n" +
			"HASH: 'a # b'" + "\n" +
			"TRUNCATED: say \\" + "\n",
		scenarioType: "decode",
	},
	{
		description:    "Decode a multi-line value",
		subdescription: "A quoted value may span lines, which is how a PEM key is written in a .env file.",
		input: "" +
			"KEY=\"-----BEGIN-----" + "\n" +
			"abc" + "\n" +
			"-----END-----\"" + "\n" +
			"NEXT=1" + "\n",
		expected: "" +
			"KEY: |-" + "\n" +
			"  -----BEGIN-----" + "\n" +
			"  abc" + "\n" +
			"  -----END-----" + "\n" +
			"NEXT: \"1\"" + "\n",
		scenarioType: "decode",
	},
	{
		description:    "Decode the awkward lines",
		subdescription: "CRLF endings, a byte order mark, an `#` inside an unquoted value (which starts a comment, as in Node), an `=` in the value, and an empty value.",
		input: "\uFEFF" +
			"WINDOWS=crlf\r\n" +
			"HASH=a#b\r\n" +
			"EQUALS=a=b\r\n" +
			"EMPTY=\r\n",
		expected: "" +
			"WINDOWS: crlf" + "\n" +
			"HASH: a" + "\n" +
			"EQUALS: a=b" + "\n" +
			"EMPTY: \"\"" + "\n",
		scenarioType: "decode",
	},
	{
		description:    "Encode a dotenv file",
		subdescription: "Nested keys are joined with an underscore and sequences are indexed, as an environment variable cannot be nested.",
		input: "" +
			"# comment" + "\n" +
			"name: Mike Wazowski" + "\n" +
			"eyes:" + "\n" +
			"  color: turquoise" + "\n" +
			"  number: 1" + "\n" +
			"friends:" + "\n" +
			"  - James P. Sullivan" + "\n" +
			"  - Celia Mae",
		expected: "" +
			"name='Mike Wazowski'" + "\n" +
			"eyes_color=turquoise" + "\n" +
			"eyes_number=1" + "\n" +
			"friends_0='James P. Sullivan'" + "\n" +
			"friends_1='Celia Mae'" + "\n",
		scenarioType: "encode",
	},
	{
		description:    "Encode values that need quoting",
		subdescription: "Single quotes are used where every reader takes them literally; double quotes only for a value that spans lines, where `\\n` is the escape.",
		input: "" +
			"plain: /usr/local/bin" + "\n" +
			"spaced: two words" + "\n" +
			"hash: \"a # b\"" + "\n" +
			"dollar: $HOME" + "\n" +
			"backslash: \"C:\\\\temp\\\\x\"" + "\n" +
			"quoted: 'say \"hi\"'" + "\n" +
			"apostrophe: \"it's fine\"" + "\n" +
			"multiline: |-" + "\n" +
			"  one" + "\n" +
			"  two" + "\n" +
			"empty: \"\"" + "\n",
		expected: "" +
			"plain=/usr/local/bin" + "\n" +
			"spaced='two words'" + "\n" +
			"hash='a # b'" + "\n" +
			"dollar='$HOME'" + "\n" +
			"backslash='C:\\temp\\x'" + "\n" +
			"quoted='say \"hi\"'" + "\n" +
			"apostrophe=\"it's fine\"" + "\n" +
			"multiline=\"one\\ntwo\"" + "\n" +
			"empty=\"\"" + "\n",
		scenarioType: "encode",
	},
	{
		description:    "Round trip a dotenv file",
		subdescription: "Decoding and re-encoding keeps every value, dropping only the comments.",
		input: "" +
			"# comment" + "\n" +
			"NAME=\"Mike Wazowski\"" + "\n" +
			"PATH=/usr/local/bin" + "\n" +
			"GREETING=\"hello world\"" + "\n" +
			"QUOTED='say \"hi\"'" + "\n",
		expected: "" +
			"NAME='Mike Wazowski'" + "\n" +
			"PATH=/usr/local/bin" + "\n" +
			"GREETING='hello world'" + "\n" +
			"QUOTED='say \"hi\"'" + "\n",
		scenarioType: "roundtrip",
	},
	{
		description:   "Decode a line without a value",
		input:         "NAME",
		expectedError: "bad file 'sample.yml': line 1: expected KEY=value, got \"NAME\"",
		scenarioType:  "decode-error",
	},
	{
		description:   "Decode an unterminated quote",
		input:         `NAME="Mike`,
		expectedError: `bad file 'sample.yml': line 1: unterminated double-quoted value "\"Mike"`,
		scenarioType:  "decode-error",
	},
	{
		description:   "Encode a value dotenv cannot express",
		input:         "mixed: >-\n  it's \"quoted\"\n",
		expectedError: `cannot encode mixed: value "it's \"quoted\"" has both quote characters, which dotenv cannot express`,
		scenarioType:  "encode-error",
	},
}

func testDotenvScenario(t *testing.T, s formatScenario) {
	switch s.scenarioType {
	case "decode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewDotenvDecoder(), NewYamlEncoder(ConfiguredYamlPreferences)), s.description)
	case "encode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewDotenvEncoder()), s.description)
	case "roundtrip":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewDotenvDecoder(), NewDotenvEncoder()), s.description)
	case "decode-error", "encode-error":
		result, err := processFormatScenario(s, NewDotenvDecoder(), NewYamlEncoder(ConfiguredYamlPreferences))
		if s.scenarioType == "encode-error" {
			result, err = processFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewDotenvEncoder())
		}
		if err == nil {
			t.Errorf("Expected error '%v' but it worked: %v", s.expectedError, result)
		} else {
			test.AssertResultComplexWithContext(t, s.expectedError, err.Error(), s.description)
		}
	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentDotenvScenario(_ *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(formatScenario)
	if s.skipDoc {
		return
	}
	switch s.scenarioType {
	case "encode":
		documentEncodeDotenvScenario(w, s)
	case "roundtrip":
		documentRoundtripDotenvScenario(w, s)
	case "decode", "decode-error":
		documentDecodeDotenvScenario(w, s)
	case "encode-error":
		documentEncodeErrorDotenvScenario(w, s)
	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func documentDecodeDotenvScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.env file of:\n")
	writeOrPanic(w, fmt.Sprintf("```env\n%v\n```\n", s.input))

	if s.scenarioType == "decode-error" {
		writeOrPanic(w, "then an error is expected:\n")
		writeOrPanic(w, fmt.Sprintf("```\n%v\n```\n\n", s.expectedError))
		return
	}

	writeOrPanic(w, "then\n")
	expression := s.expression
	if expression == "" {
		expression = "."
	}
	writeOrPanic(w, fmt.Sprintf("```bash\nyq -p=dotenv '%v' sample.env\n```\n", expression))
	writeOrPanic(w, "will output\n")
	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", mustProcessFormatScenario(s, NewDotenvDecoder(), NewYamlEncoder(ConfiguredYamlPreferences))))
}

func documentEncodeDotenvScenario(w *bufio.Writer, s formatScenario) {
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
	writeOrPanic(w, fmt.Sprintf("```bash\nyq -o=dotenv '%v' sample.yml\n```\n", expression))
	writeOrPanic(w, "will output\n")
	writeOrPanic(w, fmt.Sprintf("```env\n%v```\n\n", mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewDotenvEncoder())))
}

func documentEncodeErrorDotenvScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	writeOrPanic(w, "Given a sample.yml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```yaml\n%v\n```\n", s.input))

	writeOrPanic(w, "then an error is expected:\n")
	writeOrPanic(w, fmt.Sprintf("```\n%v\n```\n\n", s.expectedError))
}

func documentRoundtripDotenvScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.env file of:\n")
	writeOrPanic(w, fmt.Sprintf("```env\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	writeOrPanic(w, "```bash\nyq -p=dotenv -o=dotenv '.' sample.env\n```\n")
	writeOrPanic(w, "will output\n")
	writeOrPanic(w, fmt.Sprintf("```env\n%v```\n\n", mustProcessFormatScenario(s, NewDotenvDecoder(), NewDotenvEncoder())))
}

func TestDotenvDecoderInitResetsFinished(t *testing.T) {
	decoder := NewDotenvDecoder()
	firstDocuments, err := readDocuments(strings.NewReader("A=1\n"), "first.env", 0, decoder)
	if err != nil {
		t.Fatal(err)
	}
	test.AssertResult(t, 1, firstDocuments.Len())

	secondDocuments, err := readDocuments(strings.NewReader("B=2\n"), "second.env", 1, decoder)
	if err != nil {
		t.Fatal(err)
	}
	test.AssertResult(t, 1, secondDocuments.Len())
}

func TestDotenvScenarios(t *testing.T) {
	for _, tt := range dotenvScenarios {
		testDotenvScenario(t, tt)
	}

	genericScenarios := make([]interface{}, 0, len(dotenvScenarios))
	for _, s := range dotenvScenarios {
		genericScenarios = append(genericScenarios, s)
	}
	documentScenarios(t, "usage", "dotenv", genericScenarios, documentDotenvScenario)
}
