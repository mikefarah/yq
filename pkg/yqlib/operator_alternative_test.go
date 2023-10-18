package yqlib

import (
	"testing"
)

var alternativeOperatorScenarios = []expressionScenario{
	{
		// to match jq - we do not use a readonly clone context on the LHS.
		skipDoc:    true,
		expression: `.b // .c`,
		document:   `a: bridge`,
		expected: []string{
			"D0, P[c], (!!null)::null\n",
		},
	},
	{
		skipDoc:    true,
		expression: `(.b // "hello") as $x | .`,
		document:   `a: bridge`,
		expected: []string{
			"D0, P[], (!!map)::a: bridge\n",
		},
	},
	{
		skipDoc:    true,
		expression: `.a // .b`,
		document:   `a: 2`,
		expected: []string{
			"D0, P[a], (!!int)::2\n",
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
	{
		description:    "Update or create - entity exists",
		subdescription: "This initialises `a` if it's not present",
		expression:     "(.a // (.a = 0)) += 1",
		document:       `a: 1`,
		expected: []string{
			"D0, P[], (!!map)::a: 2\n",
		},
	},
	{
		description:    "Update or create - entity does not exist",
		subdescription: "This initialises `a` if it's not present",
		expression:     "(.a // (.a = 0)) += 1",
		document:       `b: camel`,
		expected: []string{
			"D0, P[], (!!map)::b: camel\na: 1\n",
		},
	},
}

func TestAlternativeOperatorScenarios(t *testing.T) {
	for _, tt := range alternativeOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "alternative-default-value", alternativeOperatorScenarios)
}
