package yqlib

import (
	"testing"
)

var collectOperatorScenarios = []expressionScenario{
	{
		document:   ``,
		expression: `[]`,
		expected:   []string{},
	},
	{
		document:   ``,
		expression: `["cat"]`,
		expected: []string{
			"D0, P[], (!!seq)::- cat\n",
		},
	}, {
		document:   ``,
		expression: `[true]`,
		expected: []string{
			"D0, P[], (!!seq)::- true\n",
		},
	}, {
		document:   ``,
		expression: `["cat", "dog"]`,
		expected: []string{
			"D0, P[], (!!seq)::- cat\n- dog\n",
		},
	}, {
		document:   ``,
		expression: `1 | collect`,
		expected: []string{
			"D0, P[], (!!seq)::- 1\n",
		},
	}, {
		document:   `[1,2,3]`,
		expression: `[.[]]`,
		expected: []string{
			"D0, P[], (!!seq)::- 1\n- 2\n- 3\n",
		},
	}, {
		document:   `a: {b: [1,2,3]}`,
		expression: `[.a.b[]]`,
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
