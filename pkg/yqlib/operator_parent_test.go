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
		description: "Parent of nested matches",
		document:    `{a: {fruit: apple, name: bob}, b: {fruit: banana, name: sam}}`,
		expression:  `.. | select(. == "banana") | parent`,
		expected: []string{
			"D0, P[b], (!!map)::{fruit: banana, name: sam}\n",
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
