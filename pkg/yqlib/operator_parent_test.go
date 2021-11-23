package yqlib

import (
	"testing"
)

var parentOperatorScenarios = []expressionScenario{
	{
		description: "Simple example",
		document:    `a: {nested: cat}`,
		expression:  `.a.nested | parent`,
		expected: []string{
			"D0, P[a], (!!map)::{nested: cat}\n",
		},
	},
	{
		description: "Show parent",
		document:    `{a: {fruit: apple}, b: {fruit: banana}}`,
		expression:  `.. | select(. == "banana") | parent`,
		expected: []string{
			"D0, P[b], (!!map)::{fruit: banana}\n",
		},
	},
	{
		description: "No parent",
		document:    `{}`,
		expression:  `parent`,
		expected:    []string{},
	},
}

func TestParentOperatorScenarios(t *testing.T) {
	for _, tt := range parentOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "parent", parentOperatorScenarios)
}
