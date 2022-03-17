package yqlib

import (
	"testing"
)

var compareOperatorScenarios = []expressionScenario{
	// both null
	{
		description: "Both sides are null: > is false",
		expression:  ".a > .b",
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		expression: ".a < .b",
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		description: "Both sides are null: >= is true",
		expression:  ".a >= .b",
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		expression: ".a <= .b",
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},

	// one null
	{
		description: "One side is null: > is false",
		document:    `a: 5`,
		expression:  ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: 5`,
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		description: "One side is null: >= is false",
		document:    `a: 5`,
		expression:  ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: 5`,
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: 5`,
		expression: ".b <= .a",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: 5`,
		expression: ".b < .a",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},

	// ints, not equal
	{
		description: "Compare integers (>)",
		document:    "a: 5\nb: 4",
		expression:  ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5\nb: 4",
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		description: "Compare integers (>=)",
		document:    "a: 5\nb: 4",
		expression:  ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5\nb: 4",
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},

	// ints, equal
	{
		description: "Compare equal numbers",
		document:    "a: 5\nb: 5",
		expression:  ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5\nb: 5",
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		description: "Compare equal numbers (>=)",
		document:    "a: 5\nb: 5",
		expression:  ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5\nb: 5",
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},

	// floats, not equal
	{
		skipDoc:    true,
		document:   "a: 5.2\nb: 4.1",
		expression: ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5.2\nb: 4.1",
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5.2\nb: 4.1",
		expression: ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5.5\nb: 4.1",
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},

	// floats, equal
	{
		skipDoc:    true,
		document:   "a: 5.5\nb: 5.5",
		expression: ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5.5\nb: 5.5",
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5.1\nb: 5.1",
		expression: ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5.1\nb: 5.1",
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},

	// strings, not equal
	// strings, equal

	// datetime, not equal
	// datetime, equal
}

func TestCompareOperatorScenarios(t *testing.T) {
	for _, tt := range compareOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "compare", compareOperatorScenarios)
}
