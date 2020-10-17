package treeops

import (
	"testing"
)

var unionOperatorScenarios = []expressionScenario{
	{
		document:   `{}`,
		expression: `"cat", "dog"`,
		expected: []string{
			"D0, P[], (!!str)::cat\n",
			"D0, P[], (!!str)::dog\n",
		},
	}, {
		document:   `{a: frog}`,
		expression: `1, true, "cat", .a`,
		expected: []string{
			"D0, P[], (!!int)::1\n",
			"D0, P[], (!!bool)::true\n",
			"D0, P[], (!!str)::cat\n",
			"D0, P[a], (!!str)::frog\n",
		},
	},
}

func TestUnionOperatorScenarios(t *testing.T) {
	for _, tt := range unionOperatorScenarios {
		testScenario(t, &tt)
	}
}
