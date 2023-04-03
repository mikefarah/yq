package yqlib

import (
	"bufio"
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

const propertiesWithCommentsOnMap = `this.thing = hi hi
# important notes
# about this value
this.value = cool
`

const expectedPropertiesWithCommentsOnMapProps = `this.thing = hi hi

# important notes
# about this value
this.value = cool
`

const expectedPropertiesWithCommentsOnMapYaml = `this:
  thing: hi hi
  # important notes
  # about this value
  value: cool
`

const propertiesWithCommentInArray = `
this.array.0 = cat
# important notes
# about dogs
this.array.1 = dog
`

const expectedPropertiesWithCommentInArrayProps = `this.array.0 = cat

# important notes
# about dogs
this.array.1 = dog
`

const expectedPropertiesWithCommentInArrayYaml = `this:
  array:
    - cat
    # important notes
    # about dogs
    - dog
`

const samplePropertiesYaml = `# block comments come through
person: # neither do comments on maps
    name: Mike Wazowski # comments on values appear
    pets: 
    - cat # comments on array values appear
    food: [pizza] # comments on arrays do not
emptyArray: []
emptyMap: []
`

const expectedPropertiesUnwrapped = `# block comments come through
# comments on values appear
person.name = Mike Wazowski

# comments on array values appear
person.pets.0 = cat
person.food.0 = pizza
`

const expectedPropertiesWrapped = `# block comments come through
# comments on values appear
person.name = "Mike Wazowski"

# comments on array values appear
person.pets.0 = cat
person.food.0 = pizza
`

const expectedUpdatedProperties = `# block comments come through
# comments on values appear
person.name = Mike Wazowski

# comments on array values appear
person.pets.0 = dog
person.food.0 = pizza
`

const expectedDecodedYaml = `person:
  # block comments come through
  # comments on values appear
  name: Mike Wazowski
  pets:
    # comments on array values appear
    - cat
  food:
    - pizza
`

const expectedPropertiesNoComments = `person.name = Mike Wazowski
person.pets.0 = cat
person.food.0 = pizza
`

const expectedPropertiesWithEmptyMapsAndArrays = `# block comments come through
# comments on values appear
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
	},
	{
		description:    "Encode properties: scalar encapsulation",
		subdescription: "Note that string values with blank characters in them are encapsulated with double quotes",
		input:          samplePropertiesYaml,
		expected:       expectedPropertiesWrapped,
		scenarioType:   "encode-wrapped",
	},
	{
		description: "Encode properties: no comments",
		input:       samplePropertiesYaml,
		expected:    expectedPropertiesNoComments,
		expression:  `... comments = ""`,
	},
	{
		description:    "Encode properties: include empty maps and arrays",
		subdescription: "Use a yq expression to set the empty maps and sequences to your desired value.",
		expression:     `(.. | select( (tag == "!!map" or tag =="!!seq") and length == 0)) = ""`,
		input:          samplePropertiesYaml,
		expected:       expectedPropertiesWithEmptyMapsAndArrays,
	},
	{
		description:  "Decode properties",
		input:        expectedPropertiesUnwrapped,
		expected:     expectedDecodedYaml,
		scenarioType: "decode",
	},

	{
		skipDoc:      true,
		description:  "Decode properties - keeps key information",
		input:        expectedPropertiesUnwrapped,
		expression:   ".person.name | key",
		expected:     "name",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "Decode properties - keeps parent information",
		input:        expectedPropertiesUnwrapped,
		expression:   ".person.name | parent",
		expected:     "name",
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "Decode properties - keeps path information",
		input:        expectedPropertiesUnwrapped,
		expression:   ".person.name | path",
		expected:     "name",
		scenarioType: "decode",
	},

	{
		description:    "Decode properties - array should be a map",
		subdescription: "If you have a numeric map key in your property files, use array_to_map to convert them to maps.",
		input:          `things.10 = mike`,
		expression:     `.things |= array_to_map`,
		expected:       "things:\n  10: mike\n",
		scenarioType:   "decode",
	},
	{
		description:  "does not expand automatically",
		skipDoc:      true,
		input:        "mike = ${dontExpand} this",
		expected:     "mike: ${dontExpand} this\n",
		scenarioType: "decode",
	},
	{
		description:  "print scalar",
		skipDoc:      true,
		input:        "mike = cat",
		expression:   ".mike",
		expected:     "cat\n",
		scenarioType: "roundtrip",
	},
	{
		description:  "Roundtrip",
		input:        expectedPropertiesUnwrapped,
		expression:   `.person.pets.0 = "dog"`,
		expected:     expectedUpdatedProperties,
		scenarioType: "roundtrip",
	},
	{
		skipDoc:      true,
		description:  "comments on arrays roundtrip",
		input:        propertiesWithCommentInArray,
		expected:     expectedPropertiesWithCommentInArrayProps,
		scenarioType: "roundtrip",
	},
	{
		skipDoc:      true,
		description:  "comments on arrays decode",
		input:        propertiesWithCommentInArray,
		expected:     expectedPropertiesWithCommentInArrayYaml,
		scenarioType: "decode",
	},
	{
		skipDoc:      true,
		description:  "comments on map roundtrip",
		input:        propertiesWithCommentsOnMap,
		expected:     expectedPropertiesWithCommentsOnMapProps,
		scenarioType: "roundtrip",
	},
	{
		skipDoc:      true,
		description:  "comments on map decode",
		input:        propertiesWithCommentsOnMap,
		expected:     expectedPropertiesWithCommentsOnMapYaml,
		scenarioType: "decode",
	},
	{
		description:  "Empty doc",
		skipDoc:      true,
		input:        "",
		expected:     "",
		scenarioType: "decode",
	},
}

func documentUnwrappedEncodePropertyScenario(w *bufio.Writer, s formatScenario) {
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

	writeOrPanic(w, fmt.Sprintf("```properties\n%v```\n\n", mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewPropertiesEncoder(true))))
}

func documentWrappedEncodePropertyScenario(w *bufio.Writer, s formatScenario) {
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
		writeOrPanic(w, fmt.Sprintf("```bash\nyq -o=props --unwrapScalar=false '%v' sample.yml\n```\n", expression))
	} else {
		writeOrPanic(w, "```bash\nyq -o=props --unwrapScalar=false sample.yml\n```\n")
	}
	writeOrPanic(w, "will output\n")

	writeOrPanic(w, fmt.Sprintf("```properties\n%v```\n\n", mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewPropertiesEncoder(false))))
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

	writeOrPanic(w, fmt.Sprintf("```yaml\n%v```\n\n", mustProcessFormatScenario(s, NewPropertiesDecoder(), NewYamlEncoder(s.indent, false, ConfiguredYamlPreferences))))
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

	writeOrPanic(w, fmt.Sprintf("```properties\n%v```\n\n", mustProcessFormatScenario(s, NewPropertiesDecoder(), NewPropertiesEncoder(true))))
}

func documentPropertyScenario(t *testing.T, w *bufio.Writer, i interface{}) {
	s := i.(formatScenario)
	if s.skipDoc {
		return
	}
	switch s.scenarioType {
	case "":
		documentUnwrappedEncodePropertyScenario(w, s)
	case "decode":
		documentDecodePropertyScenario(w, s)
	case "encode-wrapped":
		documentWrappedEncodePropertyScenario(w, s)
	case "roundtrip":
		documentRoundTripPropertyScenario(w, s)

	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func TestPropertyScenarios(t *testing.T) {
	for _, s := range propertyScenarios {
		switch s.scenarioType {
		case "":
			test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewPropertiesEncoder(true)), s.description)
		case "decode":
			test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewPropertiesDecoder(), NewYamlEncoder(2, false, ConfiguredYamlPreferences)), s.description)
		case "encode-wrapped":
			test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewYamlDecoder(ConfiguredYamlPreferences), NewPropertiesEncoder(false)), s.description)
		case "roundtrip":
			test.AssertResultWithContext(t, s.expected, mustProcessFormatScenario(s, NewPropertiesDecoder(), NewPropertiesEncoder(true)), s.description)

		default:
			panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
		}
	}
	genericScenarios := make([]interface{}, len(propertyScenarios))
	for i, s := range propertyScenarios {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "properties", genericScenarios, documentPropertyScenario)
}
