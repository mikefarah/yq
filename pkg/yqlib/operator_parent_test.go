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
		description:    "Get parents",
		subdescription: "Match all parents",
		document:       "{a: {b: {c: cat} } }",
		expression:     `.a.b.c | parents`,
		expected: []string{
			"D0, P[], (!!seq)::- {c: cat}\n- {b: {c: cat}}\n- {a: {b: {c: cat}}}\n",
		},
	},
	{
		description:    "Get the top (root) parent",
		subdescription: "Use negative numbers to get the top parents. You can think of this as indexing into the 'parents' array above",
		document:       "a:\n  b:\n    c: cat\n",
		expression:     `.a.b.c | parent(-1)`,
		expected: []string{
			"D0, P[], (!!map)::a:\n    b:\n        c: cat\n",
		},
	},
	{
		description:    "Root",
		subdescription: "Alias for parent(-1), returns the top level parent. This is usually the document node.",
		document:       "a:\n  b:\n    c: cat\n",
		expression:     `.a.b.c | root`,
		expected: []string{
			"D0, P[], (!!map)::a:\n    b:\n        c: cat\n",
		},
	},
	{
		description: "boundary negative",
		skipDoc:     true,
		document:    "a:\n  b:\n    c: cat\n",
		expression:  `.a.b.c | parent(-3)`,
		expected: []string{
			"D0, P[a b], (!!map)::c: cat\n",
		},
	},
	{
		description: "large negative",
		skipDoc:     true,
		document:    "a:\n  b:\n    c: cat\n",
		expression:  `.a.b.c | parent(-10)`,
		expected: []string{
			"D0, P[a b c], (!!str)::cat\n",
		},
	},
	{
		description: "parent zero",
		skipDoc:     true,
		document:    "a:\n  b:\n    c: cat\n",
		expression:  `.a.b.c | parent(0)`,
		expected: []string{
			"D0, P[a b c], (!!str)::cat\n",
		},
	},
	{
		description: "large positive",
		skipDoc:     true,
		document:    "a:\n  b:\n    c: cat\n",
		expression:  `.a.b.c | parent(10)`,
		expected:    []string{},
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
		description:    "N-th negative",
		subdescription: "Similarly, use negative numbers to index backwards from the parents array",
		document:       "a:\n  b:\n    c: cat\n",
		expression:     `.a.b.c | parent(-2)`,
		expected: []string{
			"D0, P[a], (!!map)::b:\n    c: cat\n",
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
