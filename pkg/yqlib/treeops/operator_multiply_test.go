package treeops

import (
	"testing"
)

var multiplyOperatorScenarios = []expressionScenario{
	{
		// document:   `{a: frog, b: cat}`,
		// expression: `.a * .b`,
		// expected: []string{
		// 	"D0, P[], (!!map)::{a: cat, b: cat}\n",
		// },
		// }, {
		document:   `{a: {things: great}, b: {also: me}}`,
		expression: `.a * .b`,
		expected: []string{
			"D0, P[], (!!map)::{a: {things: great, also: me}, b: {also: me}}\n",
		},
	},
}

func TestMultiplyOperatorScenarios(t *testing.T) {
	for _, tt := range multiplyOperatorScenarios {
		testScenario(t, &tt)
	}
}
