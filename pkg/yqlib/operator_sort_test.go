package yqlib

import "testing"

var sortByOperatorScenarios = []expressionScenario{
	{
		description: "Sort by string field",
		document:    "[{a: banana},{a: cat},{a: apple}]",
		expression:  `sort_by(.a)`,
		expected: []string{
			"D0, P[], (!!seq)::[{a: apple}, {a: banana}, {a: cat}]\n",
		},
	},
	{
		description: "Sort by multiple fields",
		document:    "[{a: dog},{a: cat, b: banana},{a: cat, b: apple}]",
		expression:  `sort_by(.a, .b)`,
		expected: []string{
			"D0, P[], (!!seq)::[{a: cat, b: apple}, {a: cat, b: banana}, {a: dog}]\n",
		},
	},
	{
		description: "Sort by multiple fields",
		skipDoc:     true,
		document:    "[{a: dog, b: good},{a: cat, c: things},{a: cat, b: apple}]",
		expression:  `sort_by(.a, .b)`,
		expected: []string{
			"D0, P[], (!!seq)::[{a: cat, c: things}, {a: cat, b: apple}, {a: dog, b: good}]\n",
		},
	},
	{
		description: "Sort by multiple fields",
		skipDoc:     true,
		document:    "[{a: dog, b: 0.1},{a: cat, b: 0.01},{a: cat, b: 0.001}]",
		expression:  `sort_by(.a, .b)`,
		expected: []string{
			"D0, P[], (!!seq)::[{a: cat, b: 0.001}, {a: cat, b: 0.01}, {a: dog, b: 0.1}]\n",
		},
	},
	{
		description:    "Sort descending by string field",
		subdescription: "Use sort with reverse to sort in descending order.",
		document:       "[{a: banana},{a: cat},{a: apple}]",
		expression:     `sort_by(.a) | reverse`,
		expected: []string{
			"D0, P[], (!!seq)::[{a: cat}, {a: banana}, {a: apple}]\n",
		},
	},
	{
		description: "Sort array in place",
		document:    "cool: [{a: banana},{a: cat},{a: apple}]",
		expression:  `.cool |= sort_by(.a)`,
		expected: []string{
			"D0, P[], (!!map)::cool: [{a: apple}, {a: banana}, {a: cat}]\n",
		},
	},
	{
		description:    "Sort array of objects by key",
		subdescription: "Note that you can give sort_by complex expressions, not just paths",
		document:       "cool: [{b: banana},{a: banana},{c: banana}]",
		expression:     `.cool |= sort_by(keys | .[0])`,
		expected: []string{
			"D0, P[], (!!map)::cool: [{a: banana}, {b: banana}, {c: banana}]\n",
		},
	},
	{
		description:    "Sort is stable",
		subdescription: "Note the order of the elements in unchanged when equal in sorting.",
		document:       "[{a: banana, b: 1}, {a: banana, b: 2}, {a: banana, b: 3}, {a: banana, b: 4}]",
		expression:     `sort_by(.a)`,
		expected: []string{
			"D0, P[], (!!seq)::[{a: banana, b: 1}, {a: banana, b: 2}, {a: banana, b: 3}, {a: banana, b: 4}]\n",
		},
	},
	{
		description: "Sort by numeric field",
		document:    "[{a: 10},{a: 100},{a: 1}]",
		expression:  `sort_by(.a)`,
		expected: []string{
			"D0, P[], (!!seq)::[{a: 1}, {a: 10}, {a: 100}]\n",
		},
	},
	{
		description: "Sort by custom date field",
		document:    `[{a: "12-Jun-2011"},{a: "23-Dec-2010"},{a: "10-Aug-2011"}]`,
		expression:  `with_dtf("02-Jan-2006"; sort_by(.a))`,
		expected: []string{
			"D0, P[], (!!seq)::[{a: \"23-Dec-2010\"}, {a: \"12-Jun-2011\"}, {a: \"10-Aug-2011\"}]\n",
		},
	},
	{
		skipDoc:    true,
		document:   "[{a: 1.1},{a: 1.001},{a: 1.01}]",
		expression: `sort_by(.a)`,
		expected: []string{
			"D0, P[], (!!seq)::[{a: 1.001}, {a: 1.01}, {a: 1.1}]\n",
		},
	},
	{
		description: "Sort, nulls come first",
		document:    "[8,3,null,6, true, false, cat]",
		expression:  `sort`,
		expected: []string{
			"D0, P[], (!!seq)::[null, false, true, 3, 6, 8, cat]\n",
		},
	},
	{
		skipDoc:     true,
		description: "false before true",
		document:    "[{a: false, b: 1}, {a: true, b: 2}, {a: false, b: 3}]",
		expression:  `sort_by(.a)`,
		expected: []string{
			"D0, P[], (!!seq)::[{a: false, b: 1}, {a: false, b: 3}, {a: true, b: 2}]\n",
		},
	},
	{
		skipDoc:     true,
		description: "head comment",
		document:    "# abc\n- def\n# ghi",
		expression:  `sort`,
		expected: []string{
			"D0, P[], (!!seq)::# abc\n- def\n# ghi\n",
		},
	},
}

func TestSortByOperatorScenarios(t *testing.T) {
	for _, tt := range sortByOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "sort", sortByOperatorScenarios)
}
