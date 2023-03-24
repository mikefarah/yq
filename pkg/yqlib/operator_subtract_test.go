package yqlib

import (
	"testing"
)

var subtractOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		document:   `{}`,
		expression: "(.a - .b) as $x | .",
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
	{
		description:    "Date subtraction",
		subdescription: "You can subtract durations from dates. Assumes RFC3339 date time format, see [date-time operators](https://mikefarah.gitbook.io/yq/operators/date-time-operators) for more information.",
		document:       `a: 2021-01-01T03:10:00Z`,
		expression:     `.a -= "3h10m"`,
		expected: []string{
			"D0, P[], (doc)::a: 2021-01-01T00:00:00Z\n",
		},
	},
	{
		description: "Date subtraction - only date",
		skipDoc:     true,
		document:    `a: 2021-01-01`,
		expression:  `.a -= "24h"`,
		expected: []string{
			"D0, P[], (doc)::a: 2020-12-31T00:00:00Z\n",
		},
	},
	{
		description:    "Date subtraction - custom format",
		subdescription: "Use with_dtf to specify your datetime format. See [date-time operators](https://mikefarah.gitbook.io/yq/operators/date-time-operators) for more information.",
		document:       `a: Saturday, 15-Dec-01 at 6:00AM GMT`,
		expression:     `with_dtf("Monday, 02-Jan-06 at 3:04PM MST", .a -= "3h1m")`,
		expected: []string{
			"D0, P[], (doc)::a: Saturday, 15-Dec-01 at 2:59AM GMT\n",
		},
	},
	{
		skipDoc:        true,
		description:    "Date subtraction - custom format",
		subdescription: "You can subtract durations from dates. See [date-time operators](https://mikefarah.gitbook.io/yq/operators/date-time-operators) for more information.",
		document:       `a: !cat Saturday, 15-Dec-01 at 6:00AM GMT`,
		expression:     `with_dtf("Monday, 02-Jan-06 at 3:04PM MST", .a -= "3h1m")`,
		expected: []string{
			"D0, P[], (doc)::a: !cat Saturday, 15-Dec-01 at 2:59AM GMT\n",
		},
	},
	{
		description:    "Custom types: that are really numbers",
		subdescription: "When custom tags are encountered, yq will try to decode the underlying type.",
		document:       "a: !horse 2\nb: !goat 1",
		expression:     `.a -= .b`,
		expected: []string{
			"D0, P[], (doc)::a: !horse 1\nb: !goat 1\n",
		},
	},
	{
		skipDoc:        true,
		description:    "Custom types: that are really floats",
		subdescription: "When custom tags are encountered, yq will try to decode the underlying type.",
		document:       "a: !horse 2.5\nb: !goat 1.5",
		expression:     `.a - .b`,
		expected: []string{
			"D0, P[a], (!horse)::1\n",
		},
	},
	{
		skipDoc:     true,
		description: "Custom types: that are really maps",
		document:    `[!horse {a: b, c: d}, !goat {a: b}]`,
		expression:  `. - [{"c": "d", "a": "b"}]`,
		expected: []string{
			"D0, P[], (!!seq)::[!goat {a: b}]\n",
		},
	},
}

func TestSubtractOperatorScenarios(t *testing.T) {
	for _, tt := range subtractOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "subtract", subtractOperatorScenarios)
}
