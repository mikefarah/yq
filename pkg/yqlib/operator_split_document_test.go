package yqlib

import (
	"testing"
)

var splitDocOperatorScenarios = []expressionScenario{
	{
		description: "Split empty",
		document:    ``,
		expression:  `splitDoc`,
		expected: []string{
			"D0, P[], (!!null)::\n",
		},
	},
	{
		description: "Split array",
		document:    `[{a: cat}, {b: dog}]`,
		expression:  `.[] | splitDoc`,
		expected: []string{
			"D0, P[0], (!!map)::{a: cat}\n",
			"D1, P[1], (!!map)::{b: dog}\n",
		},
	},
}

func TestSplitDocOperatorScenarios(t *testing.T) {
	for _, tt := range splitDocOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Split into Documents", splitDocOperatorScenarios)
}
