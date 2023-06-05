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
			"D0, P[], (!!map)::{a: cat, b: dog, c: same}\n",
		},
	},
	{
		skipDoc:     true,
		description: "Don't pass readonly context",
		expression:  `(3 + 4) | ({} | .b = "dog")`,
		expected: []string{
			"D0, P[], (!!map)::b: dog\n",
		},
	},
}

func TestPipeOperatorScenarios(t *testing.T) {
	for _, tt := range pipeOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "pipe", pipeOperatorScenarios)
}
