package yqlib

import (
	"testing"
)

var moduloOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		document:   `[{a: 2.5, b: 2}, {a: 2, b: 0.75}]`,
		expression: ".[] | .a % .b",
		expected: []string{
			"D0, P[0 a], (!!float)::0.5\n",
			"D0, P[1 a], (!!float)::0.5\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{}`,
		expression: "(.a / .b) as $x | .",
		expected: []string{
			"D0, P[], (doc)::{}\n",
		},
	},
	{
		description:    "Number modulo - int",
		subdescription: "If the lhs and rhs are ints then the expression will be calculated with ints.",
		document:       `{a: 13, b: 2}`,
		expression:     `.a = .a % .b`,
		expected: []string{
			"D0, P[], (doc)::{a: 1, b: 2}\n",
		},
	},
	{
		description:    "Number modulo - float",
		subdescription: "If the lhs or rhs are floats then the expression will be calculated with floats.",
		document:       `{a: 12, b: 2.5}`,
		expression:     `.a = .a % .b`,
		expected: []string{
			"D0, P[], (doc)::{a: !!float 2, b: 2.5}\n",
		},
	},
	{
		description:    "Number modulo - int by zero",
		subdescription: "If the lhs is an int and rhs is a 0 the result is an error.",
		document:       `{a: 1, b: 0}`,
		expression:     `.a = .a % .b`,
		expectedError:  "cannot modulo by 0",
	},
	{
		description:    "Number modulo - float by zero",
		subdescription: "If the lhs is a float and rhs is a 0 the result is NaN.",
		document:       `{a: 1.1, b: 0}`,
		expression:     `.a = .a % .b`,
		expected: []string{
			"D0, P[], (doc)::{a: !!float NaN, b: 0}\n",
		},
	},
	{
		skipDoc:     true,
		description: "Custom types: that are really numbers",
		document:    "a: !horse 333.975\nb: !goat 299.2",
		expression:  `.a = .a % .b`,
		expected: []string{
			"D0, P[], (doc)::a: !horse 34.775000000000034\nb: !goat 299.2\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 2\nb: !goat 2.3",
		expression: `.a = .a % .b`,
		expected: []string{
			"D0, P[], (doc)::a: !!float 2\nb: !goat 2.3\n",
		},
	},
	{
		skipDoc:     true,
		description: "Keep anchors",
		document:    "a: &horse [1]",
		expression:  `.a[1] = .a[0] % 2`,
		expected: []string{
			"D0, P[], (doc)::a: &horse [1, 1]\n",
		},
	},
	{
		skipDoc:       true,
		description:   "Modulo int by string",
		document:      "a: 123\nb: '2'",
		expression:    `.a % .b`,
		expectedError: "!!int cannot modulo by !!str",
	},
	{
		skipDoc:       true,
		description:   "Modulo string by int",
		document:      "a: 2\nb: '123'",
		expression:    `.b % .a`,
		expectedError: "!!str cannot modulo by !!int",
	},
	{
		skipDoc:       true,
		description:   "Modulo map by int",
		document:      "a: {\"a\":1}\nb: 2",
		expression:    `.a % .b`,
		expectedError: "!!map (a) cannot modulo by !!int (b)",
	},
	{
		skipDoc:       true,
		description:   "Modulo array by str",
		document:      "a: [1,2]\nb: '2'",
		expression:    `.a % .b`,
		expectedError: "!!seq (a) cannot modulo by !!str (b)",
	},
}

func TestModuloOperatorScenarios(t *testing.T) {
	for _, tt := range moduloOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "modulo", moduloOperatorScenarios)
}
