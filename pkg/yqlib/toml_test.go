package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var sampleTable = `
var = "x"

[owner.contact]
name = "Tom Preston-Werner"
age = 36
`

var sampleTableExpected = `var: x
owner:
  contact:
    name: Tom Preston-Werner
    age: 36
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
		description:  "Simple",
		input:        "A = \"hello\"\nB = 12\n",
		expected:     "A: hello\nB: 12\n",
		scenarioType: "decode",
	},
	{
		description:  "Deep paths",
		input:        "person.name = \"hello\"\nperson.address = \"12 cat st\"\n",
		expected:     "person:\n  name: hello\n  address: 12 cat st\n",
		scenarioType: "decode",
	},
	{
		description:  "Simpl nested",
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
		skipDoc:      true,
		description:  "inline table",
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
		skipDoc:      true,
		input:        sampleArrayTable,
		expected:     sampleArrayTableExpected,
		scenarioType: "decode",
	},
	{
		description:  "example with header",
		input:        sampleWithHeader,
		expected:     expectedSampleWithHeader,
		scenarioType: "decode",
	},
}

func testTomlScenario(t *testing.T, s formatScenario) {
	switch s.scenarioType {
	case "", "decode":
		test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewTomlDecoder(), NewYamlEncoder(2, false, ConfiguredYamlPreferences)), s.description)
	case "decode-error":
		result, err := processFormatScenario(s, NewTomlDecoder(), NewYamlEncoder(2, false, ConfiguredYamlPreferences))
		if err == nil {
			t.Errorf("Expected error '%v' but it worked: %v", s.expectedError, result)
		} else {
			test.AssertResultComplexWithContext(t, s.expectedError, err.Error(), s.description)
		}
	}
}

func TestTomlScenarios(t *testing.T) {
	for _, tt := range tomlScenarios {
		testTomlScenario(t, tt)
	}
	// genericScenarios := make([]interface{}, len(xmlScenarios))
	// for i, s := range xmlScenarios {
	// 	genericScenarios[i] = s
	// }
	// documentScenarios(t, "usage", "xml", genericScenarios, documentXMLScenario)
}
