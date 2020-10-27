package treeops

import (
	"testing"
)

var recursiveDescentOperatorScenarios = []expressionScenario{
	{
		document:   `cat`,
		expression: `..`,
		expected: []string{
			"D0, P[], (!!str)::cat\n",
		},
	},
	{
		document:   `{a: frog}`,
		expression: `..`,
		expected: []string{
			"D0, P[], (!!map)::{a: frog}\n",
			"D0, P[a], (!!str)::frog\n",
		},
	},
	{
		document:   `{a: {b: apple}}`,
		expression: `..`,
		expected: []string{
			"D0, P[], (!!map)::{a: {b: apple}}\n",
			"D0, P[a], (!!map)::{b: apple}\n",
			"D0, P[a b], (!!str)::apple\n",
		},
	},
	{
		document:   `[1,2,3]`,
		expression: `..`,
		expected: []string{
			"D0, P[], (!!seq)::[1, 2, 3]\n",
			"D0, P[0], (!!int)::1\n",
			"D0, P[1], (!!int)::2\n",
			"D0, P[2], (!!int)::3\n",
		},
	},
	{
		document:   `[{a: cat},2,true]`,
		expression: `..`,
		expected: []string{
			"D0, P[], (!!seq)::[{a: cat}, 2, true]\n",
			"D0, P[0], (!!map)::{a: cat}\n",
			"D0, P[0 a], (!!str)::cat\n",
			"D0, P[1], (!!int)::2\n",
			"D0, P[2], (!!bool)::true\n",
		},
	},
}

func TestRecursiveDescentOperatorScenarios(t *testing.T) {
	for _, tt := range recursiveDescentOperatorScenarios {
		testScenario(t, &tt)
	}
}
