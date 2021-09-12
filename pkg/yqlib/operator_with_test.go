package yqlib

import "testing"

var withOperatorScenarios = []expressionScenario{
	{
		description: "Update and style",
		document:    `a: {deeply: {nested: value}}`,
		expression:  `with(.a.deeply.nested ; . = "newValue" | . style="single")`,
		expected: []string{
			"D0, P[], (doc)::a: {deeply: {nested: 'newValue'}}\n",
		},
	},
}

func TestWithOperatorScenarios(t *testing.T) {
	for _, tt := range withOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "With", withOperatorScenarios)
}
