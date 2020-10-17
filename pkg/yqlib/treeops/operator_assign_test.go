package treeops

import (
	"testing"
)

var assignOperatorScenarios = []expressionScenario{
	{
		document:   `{a: {b: apple}}`,
		expression: `.a.b |= "frog"`,
		expected: []string{
			"D0, P[], (!!map)::{a: {b: frog}}\n",
		},
	}, {
		document:   `{a: {b: apple}}`,
		expression: `.a.b | (. |= "frog")`,
		expected: []string{
			"D0, P[a b], (!!str)::frog\n",
		},
	}, {
		document:   `{a: {b: apple}}`,
		expression: `.a.b |= 5`,
		expected: []string{
			"D0, P[], (!!map)::{a: {b: 5}}\n",
		},
	}, {
		document:   `{a: {b: apple}}`,
		expression: `.a.b |= 3.142`,
		expected: []string{
			"D0, P[], (!!map)::{a: {b: 3.142}}\n",
		},
		// 	document:   `{}`,
		// 	expression: `["cat", "dog"]`,
		// 	expected: []string{
		// 		"D0, P[], (!!seq)::- cat\n- dog\n",
		// 	},
		// }, {
		// 	document:   `{}`,
		// 	expression: `1 | collect`,
		// 	expected: []string{
		// 		"D0, P[], (!!seq)::- 1\n",
		// 	},
		// }, {
		// 	document:   `[1,2,3]`,
		// 	expression: `[.[]]`,
		// 	expected: []string{
		// 		"D0, P[], (!!seq)::- 1\n- 2\n- 3\n",
		// 	},
		// }, {
		// 	document:   `a: {b: [1,2,3]}`,
		// 	expression: `[.a.b[]]`,
		// 	expected: []string{
		// 		"D0, P[a b], (!!seq)::- 1\n- 2\n- 3\n",
		// 	},
	},
}

func TestAssignOperatorScenarios(t *testing.T) {
	for _, tt := range assignOperatorScenarios {
		testScenario(t, &tt)
	}
}
