package yqlib

import (
	"testing"
)

var equalsOperatorScenarios = []expressionScenario{
	{
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
			"D0, P[a cat b], (!!bool)::true\n",
			"D0, P[a pat b], (!!bool)::false\n",
		},
	},
	{
		document:   ``,
		expression: `null == null`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		document:   ``,
		expression: `null == ~`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
}

func TestEqualOperatorScenarios(t *testing.T) {
	for _, tt := range equalsOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Equal Operator", equalsOperatorScenarios)
}
