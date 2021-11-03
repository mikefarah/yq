package yqlib

import (
	"testing"
)

var flattenOperatorScenarios = []expressionScenario{
	{
		description:    "Flatten",
		subdescription: "Recursively flattens all arrays",
		document:       `[1, [2], [[3]]]`,
		expression:     `flatten`,
		expected: []string{
			"D0, P[], (doc)::[1, 2, 3]\n",
		},
	},
	{
		description: "Flatten with depth of one",
		document:    `[1, [2], [[3]]]`,
		expression:  `flatten(1)`,
		expected: []string{
			"D0, P[], (doc)::[1, 2, [3]]\n",
		},
	},
	{
		description: "Flatten empty array",
		document:    `[[]]`,
		expression:  `flatten`,
		expected: []string{
			"D0, P[], (doc)::[]\n",
		},
	},
	{
		description: "Flatten array of objects",
		document:    `[{foo: bar}, [{foo: baz}]]`,
		expression:  `flatten`,
		expected: []string{
			"D0, P[], (doc)::[{foo: bar}, {foo: baz}]\n",
		},
	},
}

func TestFlattenOperatorScenarios(t *testing.T) {
	for _, tt := range flattenOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "flatten", flattenOperatorScenarios)
}
