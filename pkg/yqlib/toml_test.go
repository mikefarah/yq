package yqlib

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/mikefarah/yq/v4/test"
)

var sampleTable = `
var = "x"

[owner.contact]
name = "Tom Preston-Werner"
age = 36
`

var tableArrayBeforeOwners = `
[[owner.addresses]]
street = "first street"

[owner]
name = "Tom Preston-Werner"
`

var expectedTableArrayBeforeOwners = `owner:
  addresses:
    - street: first street
  name: Tom Preston-Werner
`

var sampleTableExpected = `var: x
owner:
  contact:
    name: Tom Preston-Werner
    age: 36
`

var doubleArrayTable = `
[[fruits]]
name = "apple"
[[fruits.varieties]]  # nested array of tables
name = "red delicious"`

var doubleArrayTableExpected = `fruits:
  - name: apple
    varieties:
      - name: red delicious
`

var doubleArrayTableMultipleEntries = `
[[fruits]]
name = "banana"
[[fruits]]
name = "apple"
[[fruits.varieties]]  # nested array of tables
name = "red delicious"`

var doubleArrayTableMultipleEntriesExpected = `fruits:
  - name: banana
  - name: apple
    varieties:
      - name: red delicious
`

var doubleArrayTableNothingAbove = `
[[fruits.varieties]]  # nested array of tables
name = "red delicious"`

var doubleArrayTableNothingAboveExpected = `fruits:
  varieties:
    - name: red delicious
`

var doubleArrayTableEmptyAbove = `
[[fruits]]
[[fruits.varieties]]  # nested array of tables
name = "red delicious"`

var doubleArrayTableEmptyAboveExpected = `fruits:
  - varieties:
      - name: red delicious
`

var emptyArrayTableThenTable = `
[[fruits]]
[animals]
[[fruits.varieties]]  # nested array of tables
name = "red delicious"`

var emptyArrayTableThenTableExpected = `fruits:
  - varieties:
      - name: red delicious
animals: {}
`

var arrayTableThenArray = `
[[rootA.kidB]]
cat = "meow"

[rootA.kidB.kidC]
dog = "bark"`

var arrayTableThenArrayExpected = `rootA:
  kidB:
    - cat: meow
      kidC:
        dog: bark
`

var sampleArrayTable = `
[owner.contact]
name = "Tom Preston-Werner"
age = 36

[[owner.addresses]]
street = "first street"
suburb = "ok"

[[owner.addresses]]
street = "second street"
suburb = "nice"
`

var sampleArrayTableExpected = `owner:
  contact:
    name: Tom Preston-Werner
    age: 36
  addresses:
    - street: first street
      suburb: ok
    - street: second street
      suburb: nice
`

var emptyTable = `
[dependencies]
`

var emptyTableExpected = "dependencies: {}\n"

var multipleEmptyTables = `
[firstEmptyTable]
[firstTableWithContent]
key = "value"
[secondEmptyTable]
[thirdEmptyTable]
[secondTableWithContent]
key = "value"
[fourthEmptyTable]
[fifthEmptyTable]
`

var expectedMultipleEmptyTables = `firstEmptyTable: {}
firstTableWithContent:
  key: value
secondEmptyTable: {}
thirdEmptyTable: {}
secondTableWithContent:
  key: value
fourthEmptyTable: {}
fifthEmptyTable: {}
`

var sampleWithHeader = `
[servers]

[servers.alpha]
ip = "10.0.0.1"
`

var expectedSampleWithHeader = `servers:
  alpha:
    ip: 10.0.0.1
`

// Roundtrip fixtures
var rtInlineTableAttr = `name = { first = "Tom", last = "Preston-Werner" }
`

var rtTableSection = `[owner.contact]
name = "Tom"
age = 36
`

var rtArrayOfTables = `[[fruits]]
name = "apple"
[[fruits.varieties]]
name = "red delicious"
`

var rtArraysAndScalars = `A = ["hello", ["world", "again"]]
B = 12
`

var rtSimple = `A = "hello"
B = 12
`

var rtDeepPaths = `[person]
name = "hello"
address = "12 cat st"
`

var rtEmptyArray = `A = []
`

var rtSampleTable = `var = "x"

[owner.contact]
name = "Tom Preston-Werner"
age = 36
`

var rtEmptyTable = `[dependencies]
`

var rtComments = `# This is a comment
A = "hello"  # inline comment
B = 12

# Table comment
[person]
name = "Tom"  # name comment
`

// Reproduce bug for https://github.com/mikefarah/yq/issues/2588
// Bug: standalone comments inside a table cause subsequent key-values to be assigned at root.
var issue2588RustToolchainWithComments = `[owner]
# comment
name = "Tomer"
`

var tableWithComment = `[owner]
# comment
[things]
`

var sampleFromWeb = `# This is a TOML document
title = "TOML Example"

[owner]
name = "Tom Preston-Werner"
dob = 1979-05-27T07:32:00-08:00

[database]
enabled = true
ports = [8000, 8001, 8002]
data = [["delta", "phi"], [3.14]]
temp_targets = { cpu = 79.5, case = 72.0 }

# [servers] yq can't do this one yet
[servers.alpha]
ip = "10.0.0.1"
role = "frontend"

[servers.beta]
ip = "10.0.0.2"
role = "backend"
`

var subArrays = `
[[array]]

[[array.subarray]]

[[array.subarray.subsubarray]]
`

var tomlTableWithComments = `[section]
the_array = [
  # comment
  "value 1",

  # comment
  "value 2",
]
`

var expectedSubArrays = `array:
  - subarray:
      - subsubarray:
          - {}
`

var tomlScenarios = []formatScenario{
	{
		skipDoc:      true,
		description:  "blank",
		input:        "",
		expected:     "",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "table array before owners",
		input:        tableArrayBeforeOwners,
		expected:     expectedTableArrayBeforeOwners,
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "datetime",
		input:        "A = 1979-05-27T07:32:00-08:00",
		expected:     "A: 1979-05-27T07:32:00-08:00\n",
		scenarioType: "decode",
	},
	{
		skipDoc:       true,
		description:   "blank",
		input:         `A = "hello`,
		expectedError: `bad file 'sample.yml': basic string not terminated by "`,
		scenarioType:  "decode-error",
	},
	{
		description:  "Parse: Simple",
		input:        "A = \"hello\"\nB = 12\n",
		expected:     "A: hello\nB: 12\n",
		scenarioType: "decode",
	},
	{
		description:  "Parse: Deep paths",
		input:        "person.name = \"hello\"\nperson.address = \"12 cat st\"\n",
		expected:     "person:\n  name: hello\n  address: 12 cat st\n",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "Parse: include key information",
		input:        "person.name = \"hello\"\nperson.address = \"12 cat st\"\n",
		expression:   ".person.name | key",
		expected:     "name\n",
		scenarioType: "roundtrip",
	},
	{
		skipDoc:      true,
		description:  "Parse: include parent information",
		input:        "person.name = \"hello\"\nperson.address = \"12 cat st\"\n",
		expression:   ".person.name | parent",
		expected:     "name: hello\naddress: 12 cat st\n",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "Parse: include path information",
		input:        "person.name = \"hello\"\nperson.address = \"12 cat st\"\n",
		expression:   ".person.name | path",
		expected:     "- person\n- name\n",
		scenarioType: "decode",
	},
	{
		description:  "Encode: Scalar",
		input:        "person.name = \"hello\"\nperson.address = \"12 cat st\"\n",
		expression:   ".person.name",
		expected:     "hello\n",
		scenarioType: "roundtrip",
	},
	{
		skipDoc:      true,
		input:        `A.B = "hello"`,
		expected:     "A:\n  B: hello\n",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "bool",
		input:        `A = true`,
		expected:     "A: true\n",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "bool false",
		input:        `A = false `,
		expected:     "A: false\n",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "number",
		input:        `A = 3 `,
		expected:     "A: 3\n",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "number",
		input:        `A = 0xDEADBEEF`,
		expression:   " .A += 1",
		expected:     "A: 0xDEADBEF0\n",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "float",
		input:        `A = 6.626e-34`,
		expected:     "A: 6.626e-34\n",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "empty arraY",
		input:        `A = []`,
		expected:     "A: []\n",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "array",
		input:        `A = ["hello", ["world", "again"]]`,
		expected:     "A:\n  - hello\n  - - world\n    - again\n",
		scenarioType: "decode",
	},
	{
		description:  "Parse: inline table",
		input:        `name = { first = "Tom", last = "Preston-Werner" }`,
		expected:     "name:\n  first: Tom\n  last: Preston-Werner\n",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		input:        sampleTable,
		expected:     sampleTableExpected,
		scenarioType: "decode",
	},
	{
		description:  "Parse: Array Table",
		input:        sampleArrayTable,
		expected:     sampleArrayTableExpected,
		scenarioType: "decode",
	},
	{
		description:  "Parse: Array of Array Table",
		input:        doubleArrayTable,
		expected:     doubleArrayTableExpected,
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "Parse: Array of Array Table; nothing above",
		input:        doubleArrayTableNothingAbove,
		expected:     doubleArrayTableNothingAboveExpected,
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "Parse: Array of Array Table; empty above",
		input:        doubleArrayTableEmptyAbove,
		expected:     doubleArrayTableEmptyAboveExpected,
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "Parse: Array of Array Table; multiple entries",
		input:        doubleArrayTableMultipleEntries,
		expected:     doubleArrayTableMultipleEntriesExpected,
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "Parse: Array of Array Table; then table; then array table",
		input:        emptyArrayTableThenTable,
		expected:     emptyArrayTableThenTableExpected,
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "Parse: Array of Array Table; then table",
		input:        arrayTableThenArray,
		expected:     arrayTableThenArrayExpected,
		scenarioType: "decode",
	},
	{
		description:  "Parse: Empty Table",
		input:        emptyTable,
		expected:     emptyTableExpected,
		scenarioType: "decode",
	},
	{
		description:  "Parse: with header",
		skipDoc:      true,
		input:        sampleWithHeader,
		expected:     expectedSampleWithHeader,
		scenarioType: "decode",
	},
	{
		description:  "Parse: multiple empty tables",
		skipDoc:      true,
		input:        multipleEmptyTables,
		expected:     expectedMultipleEmptyTables,
		scenarioType: "decode",
	},
	{
		description:  "subArrays",
		skipDoc:      true,
		input:        subArrays,
		expected:     expectedSubArrays,
		scenarioType: "decode",
	},
	// Roundtrip scenarios
	{
		description:  "Roundtrip: inline table attribute",
		input:        rtInlineTableAttr,
		expression:   ".",
		expected:     rtInlineTableAttr,
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: table section",
		input:        rtTableSection,
		expression:   ".",
		expected:     rtTableSection,
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: array of tables",
		input:        rtArrayOfTables,
		expression:   ".",
		expected:     rtArrayOfTables,
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: arrays and scalars",
		input:        rtArraysAndScalars,
		expression:   ".",
		expected:     rtArraysAndScalars,
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: simple",
		input:        rtSimple,
		expression:   ".",
		expected:     rtSimple,
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: deep paths",
		input:        rtDeepPaths,
		expression:   ".",
		expected:     rtDeepPaths,
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: empty array",
		input:        rtEmptyArray,
		expression:   ".",
		expected:     rtEmptyArray,
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: sample table",
		input:        rtSampleTable,
		expression:   ".",
		expected:     rtSampleTable,
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: empty table",
		input:        rtEmptyTable,
		expression:   ".",
		expected:     rtEmptyTable,
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: comments",
		input:        rtComments,
		expression:   ".",
		expected:     rtComments,
		scenarioType: "roundtrip",
	},
	{
		skipDoc:      true,
		description:  "Issue #2588: comments inside table must not flatten (.owner.name)",
		input:        issue2588RustToolchainWithComments,
		expression:   ".owner.name",
		expected:     "Tomer\n",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "Issue #2588: comments inside table must not flatten (.name)",
		input:        issue2588RustToolchainWithComments,
		expression:   ".name",
		expected:     "null\n",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		input:        issue2588RustToolchainWithComments,
		expected:     issue2588RustToolchainWithComments,
		scenarioType: "roundtrip",
	},
	{
		skipDoc:      true,
		input:        tableWithComment,
		expression:   ".owner | headComment",
		expected:     "comment\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip: sample from web",
		input:        sampleFromWeb,
		expression:   ".",
		expected:     sampleFromWeb,
		scenarioType: "roundtrip",
	},
	{
		skipDoc:      true,
		input:        tomlTableWithComments,
		expected:     tomlTableWithComments,
		scenarioType: "roundtrip",
	},
}

func testTomlScenario(t *testing.T, s formatScenario) {
	switch s.scenarioType {
	case "", "decode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewTomlDecoder(), NewYamlEncoder(ConfiguredYamlPreferences)), s.description)
	case "decode-error":
		result, err := processFormatScenario(s, NewTomlDecoder(), NewYamlEncoder(ConfiguredYamlPreferences))
		if err == nil {
			t.Errorf("Expected error '%v' but it worked: %v", s.expectedError, result)
		} else {
			test.AssertResultComplexWithContext(t, s.expectedError, err.Error(), s.description)
		}
	case "roundtrip":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewTomlDecoder(), NewTomlEncoder()), s.description)
	}
}

func documentTomlDecodeScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.toml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```toml\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	expression := s.expression
	if expression == "" {
		expression = "."
	}
	writeOrPanic(w, fmt.Sprintf("```bash\nyq -oy '%v' sample.toml\n```\n", expression))
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", mustProcessFormatScenario(s, NewTomlDecoder(), NewYamlEncoder(ConfiguredYamlPreferences))))
}

func documentTomlRoundtripScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.toml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```toml\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")
	expression := s.expression
	if expression == "" {
		expression = "."
	}
	writeOrPanic(w, fmt.Sprintf("```bash\nyq '%v' sample.toml\n```\n", expression))
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", mustProcessFormatScenario(s, NewTomlDecoder(), NewTomlEncoder())))
}

func documentTomlScenario(_ *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(formatScenario)

	if s.skipDoc {
		return
	}
	switch s.scenarioType {
	case "", "decode":
		documentTomlDecodeScenario(w, s)
	case "roundtrip":
		documentTomlRoundtripScenario(w, s)

	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func TestTomlScenarios(t *testing.T) {
	for _, tt := range tomlScenarios {
		testTomlScenario(t, tt)
	}
	genericScenarios := make([]interface{}, len(tomlScenarios))
	for i, s := range tomlScenarios {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "toml", genericScenarios, documentTomlScenario)
}

// TestTomlColourization tests that colourization correctly distinguishes
// between table section headers and inline arrays
func TestTomlColourization(t *testing.T) {
	// Save and restore color state
	oldNoColor := color.NoColor
	color.NoColor = false
	defer func() { color.NoColor = oldNoColor }()

	// Test that inline arrays are not coloured as table sections
	encoder := &tomlEncoder{prefs: TomlPreferences{ColorsEnabled: true}}

	// Create TOML with both table sections and inline arrays
	input := []byte(`[database]
enabled = true
ports = [8000, 8001, 8002]

[servers]
alpha = "test"
`)

	result := encoder.colorizeToml(input)
	resultStr := string(result)

	// The bug would cause the inline array [8000, 8001, 8002] to be
	// coloured with the section colour (Yellow + Bold) instead of being
	// left uncoloured or coloured differently.
	//
	// To test this, we check that the section colour codes appear only
	// for actual table sections, not for inline arrays.

	// Get the ANSI codes for section colour (Yellow + Bold)
	sectionColourObj := color.New(color.FgYellow, color.Bold)
	sectionColourObj.EnableColor()
	sampleSection := sectionColourObj.Sprint("[database]")

	// Extract just the ANSI codes from the sample
	// ANSI codes start with \x1b[
	var ansiStart string
	for i := 0; i < len(sampleSection); i++ {
		if sampleSection[i] == '\x1b' {
			// Find the end of the ANSI sequence (ends with 'm')
			end := i
			for end < len(sampleSection) && sampleSection[end] != 'm' {
				end++
			}
			if end < len(sampleSection) {
				ansiStart = sampleSection[i : end+1]
				break
			}
		}
	}

	// Count how many times the section colour appears in the output
	// It should appear exactly twice: once for [database] and once for [servers]
	// If it appears more times (e.g., for [8000, 8001, 8002]), that's the bug
	sectionColourCount := strings.Count(resultStr, ansiStart)

	// We expect exactly 2 occurrences (for [database] and [servers])
	// The bug would cause more occurrences (e.g., also for [8000)
	if sectionColourCount != 2 {
		t.Errorf("Expected section colour to appear exactly 2 times (for [database] and [servers]), but it appeared %d times.\nOutput: %s", sectionColourCount, resultStr)
	}
}

func TestTomlColorisationNumberBug(t *testing.T) {
	// Save and restore color state
	oldNoColor := color.NoColor
	color.NoColor = false
	defer func() { color.NoColor = oldNoColor }()

	encoder := NewTomlEncoder()
	tomlEncoder := encoder.(*tomlEncoder)

	// Test case that exposes the bug: "123-+-45" should NOT be colourised as a single number
	input := "A = 123-+-45\n"
	result := string(tomlEncoder.colorizeToml([]byte(input)))

	// The bug causes "123-+-45" to be colourised as one token
	// It should stop at "123" because the next character '-' is not valid in this position
	if strings.Contains(result, "123-+-45") {
		// Check if it's colourised as a single token (no color codes in the middle)
		idx := strings.Index(result, "123-+-45")
		// Look backwards for color code
		beforeIdx := idx - 1
		for beforeIdx >= 0 && result[beforeIdx] != '\x1b' {
			beforeIdx--
		}
		// Look forward for reset code
		afterIdx := idx + 8 // length of "123-+-45"
		hasResetAfter := false
		for afterIdx < len(result) && afterIdx < idx+20 {
			if result[afterIdx] == '\x1b' {
				hasResetAfter = true
				break
			}
			afterIdx++
		}

		if beforeIdx >= 0 && hasResetAfter {
			// The entire "123-+-45" is wrapped in color codes - this is the bug!
			t.Errorf("BUG DETECTED: '123-+-45' is incorrectly colourised as a single number")
			t.Errorf("Expected only '123' to be colourised as a number, but got the entire '123-+-45'")
			t.Logf("Full output: %q", result)
			t.Fail()
		}
	}

	// Additional test cases for the bug
	bugTests := []struct {
		name            string
		input           string
		invalidSequence string
		description     string
	}{
		{
			name:            "consecutive minuses",
			input:           "A = 123--45\n",
			invalidSequence: "123--45",
			description:     "'123--45' should not be colourised as a single number",
		},
		{
			name:            "plus in middle",
			input:           "A = 123+45\n",
			invalidSequence: "123+45",
			description:     "'123+45' should not be colourised as a single number",
		},
	}

	for _, tt := range bugTests {
		t.Run(tt.name, func(t *testing.T) {
			result := string(tomlEncoder.colorizeToml([]byte(tt.input)))
			if strings.Contains(result, tt.invalidSequence) {
				idx := strings.Index(result, tt.invalidSequence)
				beforeIdx := idx - 1
				for beforeIdx >= 0 && result[beforeIdx] != '\x1b' {
					beforeIdx--
				}
				afterIdx := idx + len(tt.invalidSequence)
				hasResetAfter := false
				for afterIdx < len(result) && afterIdx < idx+20 {
					if result[afterIdx] == '\x1b' {
						hasResetAfter = true
						break
					}
					afterIdx++
				}

				if beforeIdx >= 0 && hasResetAfter {
					t.Errorf("BUG: %s", tt.description)
					t.Logf("Full output: %q", result)
				}
			}
		})
	}

	// Test that valid scientific notation still works
	validTests := []struct {
		name  string
		input string
	}{
		{"scientific positive", "A = 1.23e+45\n"},
		{"scientific negative", "A = 6.626e-34\n"},
		{"scientific uppercase", "A = 1.23E+10\n"},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			result := tomlEncoder.colorizeToml([]byte(tt.input))
			if len(result) == 0 {
				t.Error("Expected non-empty colourised output")
			}
		})
	}
}

// Tests that the encoder handles empty path slices gracefully
func TestTomlEmptyPathPanic(t *testing.T) {
	encoder := NewTomlEncoder()
	tomlEncoder := encoder.(*tomlEncoder)

	var buf bytes.Buffer

	// Create a simple scalar node
	scalarNode := &CandidateNode{
		Kind:  ScalarNode,
		Tag:   "!!str",
		Value: "test",
	}

	// Test with empty path - this should not panic
	err := tomlEncoder.encodeTopLevelEntry(&buf, []string{}, scalarNode)
	if err == nil {
		t.Error("Expected error when encoding with empty path, got nil")
	}

}

// TestTomlStringEscapeColourization tests that string colourization correctly
// handles escape sequences, particularly escaped quotes at the end of strings
func TestTomlStringEscapeColourization(t *testing.T) {
	// Save and restore color state
	oldNoColor := color.NoColor
	color.NoColor = false
	defer func() { color.NoColor = oldNoColor }()

	encoder := NewTomlEncoder()
	tomlEncoder := encoder.(*tomlEncoder)

	testCases := []struct {
		name        string
		input       string
		description string
	}{
		{
			name:        "escaped quote at end",
			input:       `A = "test\""` + "\n",
			description: "String ending with escaped quote should be colourised correctly",
		},
		{
			name:        "escaped backslash then quote",
			input:       `A = "test\\\""` + "\n",
			description: "String with escaped backslash followed by escaped quote",
		},
		{
			name:        "escaped quote in middle",
			input:       `A = "test\"middle"` + "\n",
			description: "String with escaped quote in the middle should be colourised correctly",
		},
		{
			name:        "multiple escaped quotes",
			input:       `A = "\"test\""` + "\n",
			description: "String with escaped quotes at start and end",
		},
		{
			name:        "escaped newline",
			input:       `A = "test\n"` + "\n",
			description: "String with escaped newline should be colourised correctly",
		},
		{
			name:        "single quote with escaped single quote",
			input:       `A = 'test\''` + "\n",
			description: "Single-quoted string with escaped single quote",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// The test should not panic and should return some output
			result := tomlEncoder.colorizeToml([]byte(tt.input))
			if len(result) == 0 {
				t.Error("Expected non-empty colourised output")
			}

			// Check that the result contains the input string (with color codes)
			// At minimum, it should contain "A" and "="
			resultStr := string(result)
			if !strings.Contains(resultStr, "A") || !strings.Contains(resultStr, "=") {
				t.Errorf("Expected output to contain 'A' and '=', got: %q", resultStr)
			}
		})
	}
}

func TestTomlEncoderPrintDocumentSeparator(t *testing.T) {
	encoder := NewTomlEncoder()
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	err := encoder.PrintDocumentSeparator(writer)
	writer.Flush()

	test.AssertResult(t, nil, err)
	test.AssertResult(t, "", buf.String())
}

func TestTomlEncoderPrintLeadingContent(t *testing.T) {
	encoder := NewTomlEncoder()
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	err := encoder.PrintLeadingContent(writer, "some content")
	writer.Flush()

	test.AssertResult(t, nil, err)
	test.AssertResult(t, "", buf.String())
}

func TestTomlEncoderCanHandleAliases(t *testing.T) {
	encoder := NewTomlEncoder()
	test.AssertResult(t, false, encoder.CanHandleAliases())
}
