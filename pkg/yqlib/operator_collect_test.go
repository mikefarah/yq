package yqlib

import (
	"testing"
)

var collectOperatorScenarios = []expressionScenario{
	{
		description: "Collect empty",
		document:    ``,
		expression:  `[]`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n",
		},
	},
	{
		description: "Collect single",
		document:    ``,
		expression:  `["cat"]`,
		expected: []string{
			"D0, P[], (!!seq)::- cat\n",
		},
	}, {
		document:   ``,
		skipDoc:    true,
		expression: `[true]`,
		expected: []string{
			"D0, P[], (!!seq)::- true\n",
		},
	},
	{
		description: "Collect many",
		document:    `{a: cat, b: dog}`,
		expression:  `[.a, .b]`,
		expected: []string{
			"D0, P[], (!!seq)::- cat\n- dog\n",
		},
	},
	{
		document:   ``,
		skipDoc:    true,
		expression: `1 | collect`,
		expected: []string{
			"D0, P[], (!!seq)::- 1\n",
		},
	},
	{
		document:   `[1,2,3]`,
		skipDoc:    true,
		expression: `[.[]]`,
		expected: []string{
			"D0, P[], (!!seq)::- 1\n- 2\n- 3\n",
		},
	},
	{
		document:   `a: {b: [1,2,3]}`,
		expression: `[.a.b.[]]`,
		skipDoc:    true,
		expected: []string{
			"D0, P[a b], (!!seq)::- 1\n- 2\n- 3\n",
		},
	},
}

func TestCollectOperatorScenarios(t *testing.T) {
	for _, tt := range collectOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Collect into Array", collectOperatorScenarios)
}
