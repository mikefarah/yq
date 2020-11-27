package yqlib

import (
	"testing"
)

var addOperatorScenarios = []expressionScenario{
	{
		description: "+= test and doc",
		expression: ".a.b+= .e.f"
		expected: []string{"add .e.g to be, return top level node"}
	},
	{
		description: "Concatenate arrays",
		document:    `{a: [1,2], b: [3,4]}`,
		expression:  `.a + .b`,
		expected: []string{
			"D0, P[a], (!!seq)::[1, 2, 3, 4]\n",
		},
	},
	{
		description: "Concatenate null to array",
		document:    `{a: [1,2]}`,
		expression:  `.a + null`,
		expected: []string{
			"D0, P[a], (!!seq)::[1, 2]\n",
		},
	},
	{
		description: "Add object to array",
		document:    `{a: [1,2], c: {cat: meow}}`,
		expression:  `.a + .c`,
		expected: []string{
			"D0, P[a], (!!seq)::[1, 2, {cat: meow}]\n",
		},
	},
	{
		description: "Add string to array",
		document:    `{a: [1,2]}`,
		expression:  `.a + "hello"`,
		expected: []string{
			"D0, P[a], (!!seq)::[1, 2, hello]\n",
		},
	},
	{
		description: "Update array (append)",
		document:    `{a: [1,2], b: [3,4]}`,
		expression:  `.a = .a + .b`,
		expected: []string{
			"D0, P[], (doc)::{a: [1, 2, 3, 4], b: [3, 4]}\n",
		},
	},
}

func TestAddOperatorScenarios(t *testing.T) {
	for _, tt := range addOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Add", addOperatorScenarios)
}
