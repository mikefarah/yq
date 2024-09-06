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
		description: "Get parent attribute",
		document:    `{a: {fruit: apple, name: bob}, b: {fruit: banana, name: sam}}`,
		expression:  `.. | select(. == "banana") | parent.name`,
		expected: []string{
			"D0, P[b name], (!!str)::sam\n",
		},
	},
	{
		description:    "N-th parent",
		subdescription: "You can optionally supply the number of levels to go up for the parent, the default being 1.",
		document:       "a:\n  b:\n    c: cat\n",
		expression:     `.a.b.c | parent(2)`,
		expected: []string{
			"D0, P[a], (!!map)::b:\n    c: cat\n",
		},
	},
	{
		description: "N-th parent - another level",
		document:    "a:\n  b:\n    c: cat\n",
		expression:  `.a.b.c | parent(3)`,
		expected: []string{
			"D0, P[], (!!map)::a:\n    b:\n        c: cat\n",
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
	documentOperatorScenarios(t, "parent", parentOperatorScenarios)
}
