package yqlib

import (
	"bufio"
	"fmt"
	"testing"

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
