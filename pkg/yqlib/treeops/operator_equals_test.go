package treeops

import (
	"testing"
)

var equalsOperatorScenarios = []expressionScenario{
	{
		document:   `[cat,goat,dog]`,
		expression: `(.[] == "*at")`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	}, {
		document:   `[cat,goat,dog]`,
		expression: `.[] | (. == "*at")`,
		expected: []string{
			"D0, P[0], (!!bool)::true\n",
			"D0, P[1], (!!bool)::true\n",
			"D0, P[2], (!!bool)::false\n",
		},
	}, {
		document:   `[3, 4, 5]`,
		expression: `.[] | (. == 4)`,
		expected: []string{
			"D0, P[0], (!!bool)::false\n",
			"D0, P[1], (!!bool)::true\n",
			"D0, P[2], (!!bool)::false\n",
		},
	}, {
		document:   `a: { cat: {b: apple, c: whatever}, pat: {b: banana} }`,
		expression: `.a | (.[].b == "apple")`,
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
}

func TestEqualOperatorScenarios(t *testing.T) {
	for _, tt := range equalsOperatorScenarios {
		testScenario(t, &tt)
	}
}
