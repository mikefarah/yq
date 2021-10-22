package yqlib

import (
	"testing"
)

var encoderDecoderOperatorScenarios = []expressionScenario{
	{
		description: "Encode value as yaml string",
		document:    `{a: {cool: "thing"}}`,
		expression:  `.b = (.a | to_yaml)`,
		expected: []string{
			`D0, P[], (doc)::{a: {cool: "thing"}, b: "{cool: \"thing\"}\n"}
`,
		},
	},
	{
		description:    "Encode value as yaml string, using toyaml",
		subdescription: "Does the same thing as to_yaml, matching jq naming convention.",
		document:       `{a: {cool: "thing"}}`,
		expression:     `.b = (.a | to_yaml)`,
		expected: []string{
			`D0, P[], (doc)::{a: {cool: "thing"}, b: "{cool: \"thing\"}\n"}
`,
		},
	},
	{
		description: "Encode value as json string",
		document:    `{a: {cool: "thing"}}`,
		expression:  `.b = (.a | to_json)`,
		expected: []string{
			`D0, P[], (doc)::{a: {cool: "thing"}, b: "{\n  \"cool\": \"thing\"\n}\n"}
`,
		},
	},
	{
		description: "Encode value as props string",
		document:    `{a: {cool: "thing"}}`,
		expression:  `.b = (.a | to_props)`,
		expected: []string{
			`D0, P[], (doc)::{a: {cool: "thing"}, b: "cool = thing\n"}
`,
		},
	},
	{
		description: "Decode a yaml encoded string",
		document:    `a: "foo: bar"`,
		expression:  `.b = (.a | from_yaml)`,
		expected: []string{
			"D0, P[], (doc)::a: \"foo: bar\"\nb:\n    foo: bar\n",
		},
	},
	{
		description:           "Update an encoded yaml string",
		dontFormatInputForDoc: true,
		document:              "a: |\n  foo: bar",
		expression:            `.a |= (from_yaml | .foo = "cat" | to_yaml)`,
		expected: []string{
			"D0, P[], (doc)::a: |\n    foo: cat\n",
		},
	},
}

func TestEncoderDecoderOperatorScenarios(t *testing.T) {
	for _, tt := range encoderDecoderOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Encoder and Decoder", encoderDecoderOperatorScenarios)
}
