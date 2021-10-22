package yqlib

import (
	"testing"
)

var encoderOperatorScenarios = []expressionScenario{
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
}

func TestEncoderOperatorScenarios(t *testing.T) {
	for _, tt := range encoderOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Encoder", encoderOperatorScenarios)
}
