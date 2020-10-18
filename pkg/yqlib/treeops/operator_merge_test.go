package treeops

import (
	"testing"
)

var mergeOperatorScenarios = []expressionScenario{
	{
		document:   `{a: frog, b: cat}`,
		expression: `.a * .b`,
		expected: []string{
			"D0, P[], (!!map)::{a: cat, b: cat}\n",
		},
	}, {
		document:   `{a: {things: great}, b: {also: me}}`,
		expression: `.a * .b`,
		expected: []string{
			"D0, P[], (!!map)::{a: {also: me, things: great}, b: {also: me}}\n",
		},
	},
}

func TestMergeOperatorScenarios(t *testing.T) {
	for _, tt := range mergeOperatorScenarios {
		testScenario(t, &tt)
	}
}
