package yqlib

import (
	"testing"
)

var alternativeOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		expression: `.b // .c`,
		document:   `a: bridge`,
		expected:   []string{},
	},
	{
		skipDoc:    true,
		expression: `(.b // "hello") as $x`,
		document:   `a: bridge`,
		expected: []string{
			"D0, P[], (doc)::a: bridge\n",
		},
	},
	{
		description: "LHS is defined",
		expression:  `.a // "hello"`,
		document:    `{a: bridge}`,
		expected: []string{
			"D0, P[a], (!!str)::bridge\n",
		},
	},
	{
		expression: `select(tag == "seq") // "cat"`,
		skipDoc:    true,
		document:   `a: frog`,
		expected: []string{
			"D0, P[], (!!str)::cat\n",
		},
	},
	{
		description: "LHS is not defined",
		expression:  `.a // "hello"`,
		document:    `{}`,
		expected: []string{
			"D0, P[], (!!str)::hello\n",
		},
	},
	{
		description: "LHS is null",
		expression:  `.a // "hello"`,
		document:    `{a: ~}`,
		expected: []string{
			"D0, P[], (!!str)::hello\n",
		},
	},
	{
		description: "LHS is false",
		expression:  `.a // "hello"`,
		document:    `{a: false}`,
		expected: []string{
			"D0, P[], (!!str)::hello\n",
		},
	},
	{
		description: "RHS is an expression",
		expression:  `.a // .b`,
		document:    `{a: false, b: cat}`,
		expected: []string{
			"D0, P[b], (!!str)::cat\n",
		},
	},
	{
		skipDoc:    true,
		expression: `false // true`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
}

func TestAlternativeOperatorScenarios(t *testing.T) {
	for _, tt := range alternativeOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Alternative (Default value)", alternativeOperatorScenarios)
}
