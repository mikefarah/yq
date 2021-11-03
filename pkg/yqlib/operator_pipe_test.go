package yqlib

import (
	"testing"
)

var pipeOperatorScenarios = []expressionScenario{
	{
		description: "Simple Pipe",
		document:    `{a: {b: cat}}`,
		expression:  `.a | .b`,
		expected: []string{
			"D0, P[a b], (!!str)::cat\n",
		},
	},
	{
		description: "Multiple updates",
		document:    `{a: cow, b: sheep, c: same}`,
		expression:  `.a = "cat" | .b = "dog"`,
		expected: []string{
			"D0, P[], (doc)::{a: cat, b: dog, c: same}\n",
		},
	},
}

func TestPipeOperatorScenarios(t *testing.T) {
	for _, tt := range pipeOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "pipe", pipeOperatorScenarios)
}
