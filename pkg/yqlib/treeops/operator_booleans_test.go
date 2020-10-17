package treeops

import (
	"testing"
)

var booleanOperatorScenarios = []expressionScenario{
	{
		document:   `{}`,
		expression: `true or false`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	}, {
		document:   `{}`,
		expression: `false or false`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	}, {
		document:   `{a: true, b: false}`,
		expression: `.[] or (false, true)`,
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
			"D0, P[b], (!!bool)::false\n",
			"D0, P[b], (!!bool)::true\n",
		},
	},
}

func TestBooleanOperatorScenarios(t *testing.T) {
	for _, tt := range booleanOperatorScenarios {
		testScenario(t, &tt)
	}
}
