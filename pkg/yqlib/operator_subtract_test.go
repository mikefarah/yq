package yqlib

import (
	"testing"
)

var subtractOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		document:   `{}`,
		expression: "(.a - .b) as $x",
		expected: []string{
			"D0, P[], (doc)::{}\n",
		},
	},
	{
		description: "Array subtraction",
		expression:  `[1,2] - [2,3]`,
		expected: []string{
			"D0, P[], (!!seq)::- 1\n",
		},
	},
	{
		skipDoc:    true,
		expression: `[2,1,2,2] - [2,3]`,
		expected: []string{
			"D0, P[], (!!seq)::- 1\n",
		},
	},
	{
		description: "Array subtraction with nested array",
		expression:  `[[1], 1, 2] - [[1], 3]`,
		expected: []string{
			"D0, P[], (!!seq)::- 1\n- 2\n",
		},
	},
	{
		skipDoc:    true,
		expression: `[[1], 1, [[[2]]]] - [[1], [[[3]]]]`,
		expected: []string{
			"D0, P[], (!!seq)::- 1\n- - - - 2\n",
		},
	},
	{
		description:    "Array subtraction with nested object",
		subdescription: `Note that order of the keys does not matter`,
		document:       `[{a: b, c: d}, {a: b}]`,
		expression:     `. - [{"c": "d", "a": "b"}]`,
		expected: []string{
			"D0, P[], (!!seq)::[{a: b}]\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[{a: [1], c: d}, {a: [2], c: d}, {a: b}]`,
		expression: `. - [{"c": "d", "a": [1]}]`,
		expected: []string{
			"D0, P[], (!!seq)::[{a: [2], c: d}, {a: b}]\n",
		},
	},
	{
		description:    "Number subtraction - float",
		subdescription: "If the lhs or rhs are floats then the expression will be calculated with floats.",
		document:       `{a: 3, b: 4.5}`,
		expression:     `.a = .a - .b`,
		expected: []string{
			"D0, P[], (doc)::{a: -1.5, b: 4.5}\n",
		},
	},
	{
		description:    "Number subtraction - float",
		subdescription: "If the lhs or rhs are floats then the expression will be calculated with floats.",
		document:       `{a: 3, b: 4.5}`,
		expression:     `.a = .a - .b`,
		expected: []string{
			"D0, P[], (doc)::{a: -1.5, b: 4.5}\n",
		},
	},
	{
		description:    "Number subtraction - int",
		subdescription: "If both the lhs and rhs are ints then the expression will be calculated with ints.",
		document:       `{a: 3, b: 4}`,
		expression:     `.a = .a - .b`,
		expected: []string{
			"D0, P[], (doc)::{a: -1, b: 4}\n",
		},
	},
	{
		description: "Decrement numbers",
		document:    `{a: 3, b: 5}`,
		expression:  `.[] -= 1`,
		expected: []string{
			"D0, P[], (doc)::{a: 2, b: 4}\n",
		},
	},
}

func TestSubtractOperatorScenarios(t *testing.T) {
	for _, tt := range subtractOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "subtract", subtractOperatorScenarios)
}
