package yqlib

import (
	"testing"
)

var compareOperatorScenarios = []expressionScenario{
	// ints, not equal
	{
		description: "Compare numbers (>)",
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
		skipDoc:     true,
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
		skipDoc:     true,
		description: "Compare equal numbers (>)",
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
	{
		description:    "Compare strings",
		subdescription: "Compares strings by their bytecode.",
		document:       "a: zoo\nb: apple",
		expression:     ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: zoo\nb: apple",
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: zoo\nb: apple",
		expression: ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: zoo\nb: apple",
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},

	// strings, equal
	{
		skipDoc:    true,
		document:   "a: cat\nb: cat",
		expression: ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: cat\nb: cat",
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: cat\nb: cat",
		expression: ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: cat\nb: cat",
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},

	// datetime, not equal
	{
		description:    "Compare date times",
		subdescription: "You can compare date times. Assumes RFC3339 date time format, see [date-time operators](https://mikefarah.gitbook.io/yq/operators/date-time-operators) for more information.",
		document:       "a: 2021-01-01T03:10:00Z\nb: 2020-01-01T03:10:00Z",
		expression:     ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 2021-01-01T03:10:00Z\nb: 2020-01-01T03:10:00Z",
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 2021-01-01T03:10:00Z\nb: 2020-01-01T03:10:00Z",
		expression: ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 2021-01-01T03:10:00Z\nb: 2020-01-01T03:10:00Z",
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},

	// datetime, equal
	{
		skipDoc:    true,
		document:   "a: 2021-01-01T03:10:00Z\nb: 2021-01-01T03:10:00Z",
		expression: ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 2021-01-01T03:10:00Z\nb: 2021-01-01T03:10:00Z",
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 2021-01-01T03:10:00Z\nb: 2021-01-01T03:10:00Z",
		expression: ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 2021-01-01T03:10:00Z\nb: 2021-01-01T03:10:00Z",
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
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
}

func TestCompareOperatorScenarios(t *testing.T) {
	for _, tt := range compareOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "compare", compareOperatorScenarios)
}
