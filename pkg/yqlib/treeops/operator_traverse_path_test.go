package treeops

import (
	"testing"
)

var traversePathOperatorScenarios = []expressionScenario{
	{
		document:   `{a: {b: apple}}`,
		expression: `.a`,
		expected: []string{
			"D0, P[a], (!!map)::{b: apple}\n",
		},
	},
}

func TestTraversePathOperatorScenarios(t *testing.T) {
	for _, tt := range traversePathOperatorScenarios {
		testScenario(t, &tt)
	}
}
