package yqlib

import (
	"testing"
)

var lengthOperatorScenarios = []expressionScenario{
	{
		description:    "String length",
		subdescription: "returns length of string",
		document:       `{a: cat}`,
		expression:     `.a | length`,
		expected: []string{
			"D0, P[a], (!!int)::3\n",
		},
	},
	{
		description:    "Map length",
		subdescription: "returns number of entries",
		document:       `{a: cat, c: dog}`,
		expression:     `length`,
		expected: []string{
			"D0, P[], (!!int)::2\n",
		},
	},
	{
		description:    "Array length",
		subdescription: "returns number of elements",
		document:       `[2,4,6,8]`,
		expression:     `length`,
		expected: []string{
			"D0, P[], (!!int)::4\n",
		},
	},
}

func TestLengthOperatorScenarios(t *testing.T) {
	for _, tt := range lengthOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Length", lengthOperatorScenarios)
}
