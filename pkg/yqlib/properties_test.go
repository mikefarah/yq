package yqlib

import (
	"bufio"
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

const samplePropertiesYaml = `# block comments don't come through
person: # neither do comments on maps
    name: Mike Wazowski # comments on values appear
    pets: 
    - cat # comments on array values appear
    food: [pizza] # comments on arrays do not
emptyArray: []
emptyMap: []
`

const expectedPropertiesUnwrapped = `# comments on values appear
person.name = Mike Wazowski

# comments on array values appear
person.pets.0 = cat
person.food.0 = pizza
`

const expectedPropertiesWrapped = `# comments on values appear
person.name = "Mike Wazowski"

# comments on array values appear
person.pets.0 = cat
person.food.0 = pizza
`

const expectedUpdatedProperties = `# comments on values appear
person.name = Mike Wazowski

# comments on array values appear
person.pets.0 = dog
person.food.0 = pizza
`

const expectedDecodedYaml = `person:
  name: Mike Wazowski # comments on values appear
  pets:
    - cat # comments on array values appear
  food:
    - pizza
`

const expectedPropertiesNoComments = `person.name = Mike Wazowski
person.pets.0 = cat
person.food.0 = pizza
`

const expectedPropertiesWithEmptyMapsAndArrays = `# comments on values appear
person.name = Mike Wazowski

# comments on array values appear
person.pets.0 = cat
person.food.0 = pizza
emptyArray = 
emptyMap = 
`

var propertyScenarios = []formatScenario{
	{
		description:    "Encode properties",
		subdescription: "Note that empty arrays and maps are not encoded by default.",
		input:          samplePropertiesYaml,
		expected:       expectedPropertiesUnwrapped,
		encoder:        NewPropertiesEncoder(true),
		decoder:        NewYamlDecoder(),
	},
	{
		description:    "Encode properties: scalar encapsulation",
		subdescription: "Note that string values with blank characters in them are encapsulated with double quotes",
		input:          samplePropertiesYaml,
		expected:       expectedPropertiesWrapped,
		decoder:        NewYamlDecoder(),
		encoder:        NewPropertiesEncoder(false),
	},
	{
		description: "Encode properties: no comments",
		input:       samplePropertiesYaml,
		expected:    expectedPropertiesNoComments,
		expression:  `... comments = ""`,
		decoder:     NewYamlDecoder(),
		encoder:     NewPropertiesEncoder(true),
	},
	{
		description:    "Encode properties: include empty maps and arrays",
		subdescription: "Use a yq expression to set the empty maps and sequences to your desired value.",
		expression:     `(.. | select( (tag == "!!map" or tag =="!!seq") and length == 0)) = ""`,
		input:          samplePropertiesYaml,
		expected:       expectedPropertiesWithEmptyMapsAndArrays,
		decoder:        NewYamlDecoder(),
		encoder:        NewPropertiesEncoder(true),
	},
	{
		description:  "Decode properties",
		input:        expectedPropertiesUnwrapped,
		expected:     expectedDecodedYaml,
		scenarioType: "decode",
		decoder:      NewPropertiesDecoder(),
		encoder:      NewYamlEncoder(2, false, true, true),
	},
	{
		description:  "does not expand automatically",
		skipDoc:      true,
		input:        "mike = ${dontExpand} this",
		expected:     "mike: ${dontExpand} this\n",
		scenarioType: "decode",
		decoder:      NewPropertiesDecoder(),
		encoder:      NewYamlEncoder(2, false, true, true),
	},
	{
		description:  "Roundtrip",
		input:        expectedPropertiesUnwrapped,
		expression:   `.person.pets.0 = "dog"`,
		expected:     expectedUpdatedProperties,
		scenarioType: "roundtrip",
		decoder:      NewPropertiesDecoder(),
		encoder:      NewPropertiesEncoder(true),
	},
	{
		description:  "Empty doc",
		skipDoc:      true,
		input:        "",
		expected:     "",
		scenarioType: "decode",
		decoder:      NewPropertiesDecoder(),
		encoder:      NewYamlEncoder(2, false, true, true),
	},
}

func documentEncodePropertyScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.yml file of:\n")
	writeOrPanic(w, fmt.Sprintf("```yaml\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")

	expression := s.expression

	if expression != "" {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -o=props '%v' sample.yml\n```\n", expression))
	} else {
		writeOrPanic(w, "```bash\nyq -o=props sample.yml\n```\n")
	}
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```properties\n%v```\n\n", processFormatScenario(s, s.decoder, s.encoder)))
}

func documentDecodePropertyScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.properties file of:\n")
	writeOrPanic(w, fmt.Sprintf("```properties\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")

	expression := s.expression
	if expression != "" {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -p=props '%v' sample.properties\n```\n", expression))
	} else {
		writeOrPanic(w, "```bash\nyq -p=props sample.properties\n```\n")
	}

	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", processFormatScenario(s, s.decoder, s.encoder)))
}

func documentRoundTripPropertyScenario(w *bufio.Writer, s formatScenario) {
	writeOrPanic(w, fmt.Sprintf("## %v\n", s.description))

	if s.subdescription != "" {
		writeOrPanic(w, s.subdescription)
		writeOrPanic(w, "\n\n")
	}

	writeOrPanic(w, "Given a sample.properties file of:\n")
	writeOrPanic(w, fmt.Sprintf("```properties\n%v\n```\n", s.input))

	writeOrPanic(w, "then\n")

	expression := s.expression
	if expression != "" {
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -p=props -o=props '%v' sample.properties\n```\n", expression))
	} else {
		writeOrPanic(w, "```bash\nyq -p=props -o=props sample.properties\n```\n")
	}

	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```properties\n%v```\n\n", processFormatScenario(s, s.decoder, s.encoder)))
}

func documentPropertyScenario(t *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(formatScenario)
	if s.skipDoc {
		return
	}
	if s.scenarioType == "decode" {
		documentDecodePropertyScenario(w, s)
	} else if s.scenarioType == "roundtrip" {
		documentRoundTripPropertyScenario(w, s)
	} else {
		documentEncodePropertyScenario(w, s)
	}

}

func TestPropertyScenarios(t *testing.T) {
	for _, s := range propertyScenarios {
		if s.scenarioType == "decode" {
			test.AssertResultWithContext(t, s.expected, processFormatScenario(s, s.decoder, s.encoder), s.description)
		} else if s.scenarioType == "roundtrip" {
			test.AssertResultWithContext(t, s.expected, processFormatScenario(s, s.decoder, s.encoder), s.description)
		} else {
			test.AssertResultWithContext(t, s.expected, processFormatScenario(s, s.decoder, s.encoder), s.description)
		}
	}
	genericScenarios := make([]interface{}, len(propertyScenarios))
	for i, s := range propertyScenarios {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "properties", genericScenarios, documentPropertyScenario)
}
