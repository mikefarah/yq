package yqlib

import (
	"testing"
)

var variableOperatorScenarios = []expressionScenario{
	{
		description: "Single value variable",
		document:    `a: cat`,
		expression:  `.a as $foo | $foo`,
		expected: []string{
			"D0, P[a], (!!str)::cat\n",
		},
	},
	{
		description: "Multi value variable",
		document:    `[cat, dog]`,
		expression:  `.[] as $foo | $foo`,
		expected: []string{
			"D0, P[0], (!!str)::cat\n",
			"D0, P[1], (!!str)::dog\n",
		},
	},
}

func TestVariableOperatorScenarios(t *testing.T) {
	for _, tt := range variableOperatorScenarios {
		testScenario(t, &tt)
	}
}
