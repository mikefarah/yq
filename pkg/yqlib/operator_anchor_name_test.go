package yqlib

import (
	"testing"
)

var anchorOperatorScenarios = []expressionScenario{
	{
		description: "Get anchor",
		document:    `a: &billyBob cat`,
		expression:  `.a | anchor`,
		expected: []string{
			"D0, P[a], (!!str)::billyBob\n",
		},
	},
	{
		description: "Set anchor name",
		document:    `a: cat`,
		expression:  `.a anchor = "foobar"`,
		expected: []string{
			"D0, P[], (doc)::a: &foobar cat\n",
		},
	},
}

func TestAnchorOperatorScenarios(t *testing.T) {
	for _, tt := range anchorOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Anchor Operators", anchorOperatorScenarios)
}
