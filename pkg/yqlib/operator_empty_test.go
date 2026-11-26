package yqlib

import "testing"

var emptyOperatorScenarios = []expressionScenario{
	{
		description: "Produce no results",
		document:    `a: cat`,
		expression:  `empty`,
		expected:    []string{},
	},
	{
		description:    "Fall back to no output instead of null",
		subdescription: "Handy in shell scripts, where a missing value should result in an empty string rather than the text `null`.",
		document:       `a: cat`,
		expression:     `.b // empty`,
		expected:       []string{},
	},
	{
		description: "Fall back to no output, value present",
		document:    `a: cat`,
		expression:  `.a // empty`,
		expected: []string{
			"D0, P[a], (!!str)::cat\n",
		},
	},
	{
		description: "Empty is dropped from collected results",
		document:    `a: cat`,
		expression:  `[.a, empty, .a]`,
		expected: []string{
			"D0, P[], (!!seq)::- cat\n- cat\n",
		},
	},
	{
		skipDoc:    true,
		expression: `empty`,
		expected:   []string{},
	},
	{
		skipDoc:    true,
		document:   `a: cat`,
		expression: `.a, empty, .a`,
		expected: []string{
			"D0, P[a], (!!str)::cat\n",
			"D0, P[a], (!!str)::cat\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: cat`,
		expression: `empty | .a`,
		expected:   []string{},
	},
	{
		skipDoc:    true,
		document:   `a: cat`,
		expression: `.a | empty`,
		expected:   []string{},
	},
	{
		skipDoc:    true,
		document:   "empty: cat\n",
		expression: `.empty`,
		expected: []string{
			"D0, P[empty], (!!str)::cat\n",
		},
	},
	{
		skipDoc:    true,
		document:   "empty: cat\n",
		expression: `.["empty"]`,
		expected: []string{
			"D0, P[empty], (!!str)::cat\n",
		},
	},
}

func TestEmptyOperatorScenarios(t *testing.T) {
	for _, tt := range emptyOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "empty", emptyOperatorScenarios)
}
