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
			"D0, P[], (!!seq)::[1, 2, 3]\n",
		},
	},
	{
		description: "Flatten with depth of one",
		document:    `[1, [2], [[3]]]`,
		expression:  `flatten(1)`,
		expected: []string{
			"D0, P[], (!!seq)::[1, 2, [3]]\n",
		},
	},
	{
		description: "Flatten empty array",
		document:    `[[]]`,
		expression:  `flatten`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n",
		},
	},
	{
		description: "Flatten array of objects",
		document:    `[{foo: bar}, [{foo: baz}]]`,
		expression:  `flatten`,
		expected: []string{
			"D0, P[], (!!seq)::[{foo: bar}, {foo: baz}]\n",
		},
	},
}

func TestFlattenOperatorScenarios(t *testing.T) {
	for _, tt := range flattenOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "flatten", flattenOperatorScenarios)
}
