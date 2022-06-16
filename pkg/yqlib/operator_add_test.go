package yqlib

import (
	"testing"
)

var addOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		document:   `[{a: foo, b: bar}, {a: 1, b: 2}]`,
		expression: ".[] | .a + .b",
		expected: []string{
			"D0, P[0 a], (!!str)::foobar\n",
			"D0, P[1 a], (!!int)::3\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: key`,
		expression: `. += {"key": "b"}`,
		expected: []string{
			"D0, P[], (!!map)::a: key\nkey: b\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[[c], [b]]`,
		expression: `.[] | . += "a"`,
		expected: []string{
			"D0, P[0], (!!seq)::[c, a]\n",
			"D0, P[1], (!!seq)::[b, a]\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{}`,
		expression: "(.a + .b) as $x",
		expected: []string{
			"D0, P[], (doc)::{}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: 0`,
		expression: ".a += .b.c",
		expected: []string{
			"D0, P[], (doc)::a: 0\n",
		},
	},

	{
		description: "Concatenate arrays",
		document:    `{a: [1,2], b: [3,4]}`,
		expression:  `.a + .b`,
		expected: []string{
			"D0, P[a], (!!seq)::[1, 2, 3, 4]\n",
		},
	},
	{
		description:           "Concatenate to existing array",
		subdescription:        "Note that the styling of `a` is kept.",
		document:              "a: [1,2]\nb:\n  - 3\n  - 4",
		dontFormatInputForDoc: true,
		expression:            `.a += .b`,
		expected: []string{
			"D0, P[], (doc)::a: [1, 2, 3, 4]\nb:\n    - 3\n    - 4\n",
		},
	},
	{
		skipDoc:    true,
		expression: `[1] + ([2], [3])`,
		expected: []string{
			"D0, P[], (!!seq)::- 1\n- 2\n",
			"D0, P[], (!!seq)::- 1\n- 3\n",
		},
	},
	{
		description: "Concatenate null to array",
		document:    `{a: [1,2]}`,
		expression:  `.a + null`,
		expected: []string{
			"D0, P[a], (!!seq)::[1, 2]\n",
		},
	},
	{
		skipDoc:     true,
		description: "Concatenate to empty array",
		document:    `{a: []}`,
		expression:  `.a + "cat"`,
		expected: []string{
			"D0, P[a], (!!seq)::- cat\n",
		},
	},
	{
		description:           "Append to existing array",
		subdescription:        "Note that the styling is copied from existing array elements",
		dontFormatInputForDoc: true,
		document:              `a: ['dog']`,
		expression:            `.a += "cat"`,
		expected: []string{
			"D0, P[], (doc)::a: ['dog', 'cat']\n",
		},
	},
	{
		skipDoc:        true,
		description:    "Concatenate to existing array",
		subdescription: "does not modify original",
		document:       `{a: ['dog'], b: cat}`,
		expression:     `.a = .a + .b`,
		expected: []string{
			"D0, P[], (doc)::{a: ['dog', 'cat'], b: cat}\n",
		},
	},
	{
		skipDoc:     true,
		description: "Concatenate to empty array",
		document:    `a: []`,
		expression:  `.a += "cat"`,
		expected: []string{
			"D0, P[], (doc)::a:\n    - cat\n",
		},
	},
	{
		skipDoc:     true,
		description: "Concatenate to existing array",
		document:    `a: [dog]`,
		expression:  `.a += "cat"`,
		expected: []string{
			"D0, P[], (doc)::a: [dog, cat]\n",
		},
	},
	{
		skipDoc:     true,
		description: "Concatenate to empty object",
		document:    `{a: {}}`,
		expression:  `.a + {"b": "cat"}`,
		expected: []string{
			"D0, P[a], (!!map)::b: cat\n",
		},
	},
	{
		skipDoc:     true,
		description: "Concatenate to existing object",
		document:    `{a: {c: dog}}`,
		expression:  `.a + {"b": "cat"}`,
		expected: []string{
			"D0, P[a], (!!map)::{c: dog, b: cat}\n",
		},
	},
	{
		skipDoc:        true,
		description:    "Concatenate to existing object",
		subdescription: "matches stylig",
		document:       "a:\n  c: dog",
		expression:     `.a + {"b": "cat"}`,
		expected: []string{
			"D0, P[a], (!!map)::c: dog\nb: cat\n",
		},
	},
	{
		skipDoc:     true,
		description: "Concatenate to empty object in place",
		document:    `a: {}`,
		expression:  `.a += {"b": "cat"}`,
		expected: []string{
			"D0, P[], (doc)::a:\n    b: cat\n",
		},
	},
	{
		skipDoc:     true,
		description: "Concatenate to existing object in place",
		document:    `a: {c: dog}`,
		expression:  `.a += {"b": "cat"}`,
		expected: []string{
			"D0, P[], (doc)::a: {c: dog, b: cat}\n",
		},
	},
	{
		description: "Add new object to array",
		document:    `a: [{dog: woof}]`,
		expression:  `.a + {"cat": "meow"}`,
		expected: []string{
			"D0, P[a], (!!seq)::[{dog: woof}, {cat: meow}]\n",
		},
	},
	{
		description: "Relative append",
		document:    `a: { a1: {b: [cat]}, a2: {b: [dog]}, a3: {} }`,
		expression:  `.a[].b += ["mouse"]`,
		expected: []string{
			"D0, P[], (doc)::a: {a1: {b: [cat, mouse]}, a2: {b: [dog, mouse]}, a3: {b: [mouse]}}\n",
		},
	},
	{
		description: "String concatenation",
		document:    `{a: cat, b: meow}`,
		expression:  `.a += .b`,
		expected: []string{
			"D0, P[], (doc)::{a: catmeow, b: meow}\n",
		},
	},
	{
		description: "String concatenation - str + int",
		skipDoc:     true,
		document:    `{a: !cool cat, b: meow}`,
		expression:  `.a + 3`,
		expected: []string{
			"D0, P[a], (!cool)::cat3\n",
		},
	},
	{
		description: "String concatenation - int + str",
		skipDoc:     true,
		document:    `{a: !cool cat, b: meow}`,
		expression:  `3 + .a`,
		expected: []string{
			"D0, P[], (!cool)::3cat\n",
		},
	},
	{
		description:    "Number addition - float",
		subdescription: "If the lhs or rhs are floats then the expression will be calculated with floats.",
		document:       `{a: 3, b: 4.9}`,
		expression:     `.a = .a + .b`,
		expected: []string{
			"D0, P[], (doc)::{a: 7.9, b: 4.9}\n",
		},
	},
	{
		description:    "Number addition - int",
		subdescription: "If both the lhs and rhs are ints then the expression will be calculated with ints.",
		document:       `{a: 3, b: 4}`,
		expression:     `.a = .a + .b`,
		expected: []string{
			"D0, P[], (doc)::{a: 7, b: 4}\n",
		},
	},
	{
		description: "Increment numbers",
		document:    `{a: 3, b: 5}`,
		expression:  `.[] += 1`,
		expected: []string{
			"D0, P[], (doc)::{a: 4, b: 6}\n",
		},
	},
	{
		description:    "Date addition",
		subdescription: "You can add durations to dates. Assumes RFC3339 date time format, see [date-time operators](https://mikefarah.gitbook.io/yq/operators/date-time-operators) for more information.",
		document:       `a: 2021-01-01T00:00:00Z`,
		expression:     `.a += "3h10m"`,
		expected: []string{
			"D0, P[], (doc)::a: 2021-01-01T03:10:00Z\n",
		},
	},
	{
		description:    "Date addition - custom format",
		subdescription: "You can add durations to dates. See [date-time operators](https://mikefarah.gitbook.io/yq/operators/date-time-operators) for more information.",
		document:       `a: Saturday, 15-Dec-01 at 2:59AM GMT`,
		expression:     `with_dtf("Monday, 02-Jan-06 at 3:04PM MST", .a += "3h1m")`,
		expected: []string{
			"D0, P[], (doc)::a: Saturday, 15-Dec-01 at 6:00AM GMT\n",
		},
	},
	{
		skipDoc:        true,
		description:    "Date addition - custom format",
		subdescription: "You can add durations to dates. See [date-time operators](https://mikefarah.gitbook.io/yq/operators/date-time-operators) for more information.",
		document:       `a: !cat Saturday, 15-Dec-01 at 2:59AM GMT`,
		expression:     `with_dtf("Monday, 02-Jan-06 at 3:04PM MST", .a += "3h1m")`,
		expected: []string{
			"D0, P[], (doc)::a: !cat Saturday, 15-Dec-01 at 6:00AM GMT\n",
		},
	},
	{
		description:    "Add to null",
		subdescription: "Adding to null simply returns the rhs",
		expression:     `null + "cat"`,
		expected: []string{
			"D0, P[], (!!str)::cat\n",
		},
	},
	{
		description:    "Add maps to shallow merge",
		subdescription: "Adding objects together shallow merges them. Use `*` to deeply merge.",
		document:       "a: {thing: {name: Astuff, value: x}, a1: cool}\nb: {thing: {name: Bstuff, legs: 3}, b1: neat}",
		expression:     `.a += .b`,
		expected: []string{
			"D0, P[], (doc)::a: {thing: {name: Bstuff, legs: 3}, a1: cool, b1: neat}\nb: {thing: {name: Bstuff, legs: 3}, b1: neat}\n",
		},
	},
	{
		description:    "Custom types: that are really strings",
		subdescription: "When custom tags are encountered, yq will try to decode the underlying type.",
		document:       "a: !horse cat\nb: !goat _meow",
		expression:     `.a += .b`,
		expected: []string{
			"D0, P[], (doc)::a: !horse cat_meow\nb: !goat _meow\n",
		},
	},
	{
		description:    "Custom types: that are really numbers",
		subdescription: "When custom tags are encountered, yq will try to decode the underlying type.",
		document:       "a: !horse 1.2\nb: !goat 2.3",
		expression:     `.a += .b`,
		expected: []string{
			"D0, P[], (doc)::a: !horse 3.5\nb: !goat 2.3\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: !horse 2\nb: !goat 2.3",
		expression: `.a += .b`,
		expected: []string{
			"D0, P[], (doc)::a: !horse 4.3\nb: !goat 2.3\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 2\nb: !goat 2.3",
		expression: `.a += .b`,
		expected: []string{
			"D0, P[], (doc)::a: 4.3\nb: !goat 2.3\n",
		},
	},
	{
		skipDoc:     true,
		description: "Custom types: that are really ints",
		document:    "a: !horse 2\nb: !goat 3",
		expression:  `.a += .b`,
		expected: []string{
			"D0, P[], (doc)::a: !horse 5\nb: !goat 3\n",
		},
	},
	{
		description:    "Custom types: that are really arrays",
		skipDoc:        true,
		subdescription: "when custom tags are encountered, yq will try to decode the underlying type.",
		document:       "a: !horse [a]\nb: !goat [b]",
		expression:     `.a += .b`,
		expected: []string{
			"D0, P[], (doc)::a: !horse [a, b]\nb: !goat [b]\n",
		},
	},
}

func TestAddOperatorScenarios(t *testing.T) {
	for _, tt := range addOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "add", addOperatorScenarios)
}
