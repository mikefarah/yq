package yqlib

import (
	"testing"
)

var collectOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		document:   ``,
		expression: `.a += [0]`,
		expected: []string{
			"D0, P[], ()::a:\n    - 0\n",
		},
	},
	{
		skipDoc:    true,
		expression: `[1,2,3]`,
		expected: []string{
			"D0, P[], (!!seq)::- 1\n- 2\n- 3\n",
		},
	},
	{
		description: "Collect empty",
		document:    ``,
		expression:  `[]`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n",
		},
	},
	{
		skipDoc:  true,
		document: "{a: apple}\n---\n{b: frog}",

		expression: `[.]`,
		expected: []string{
			"D0, P[], (!!seq)::- {a: apple}\n- {b: frog}\n",
		},
	},
	{
		skipDoc:    true,
		document:   ``,
		expression: `[3]`,
		expected: []string{
			"D0, P[], (!!seq)::- 3\n",
		},
	},
	{
		description: "Collect single",
		document:    ``,
		expression:  `["cat"]`,
		expected: []string{
			"D0, P[], (!!seq)::- cat\n",
		},
	},
	{
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
		expression: `collect(1)`,
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
			"D0, P[], (!!seq)::- 1\n- 2\n- 3\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[{name: cat, thing: bor}, {name: dog}]`,
		expression: `.[] | [.name]`,
		expected: []string{
			"D0, P[0], (!!seq)::- cat\n",
			"D0, P[1], (!!seq)::- dog\n",
		},
	},
}

func TestCollectOperatorScenarios(t *testing.T) {
	for _, tt := range collectOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "collect-into-array", collectOperatorScenarios)
}
