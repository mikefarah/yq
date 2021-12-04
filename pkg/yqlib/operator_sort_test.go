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
		description: "Sort array in place",
		document:    "cool: [{a: banana},{a: cat},{a: apple}]",
		expression:  `.cool |= sort_by(.a)`,
		expected: []string{
			"D0, P[], (doc)::cool: [{a: apple}, {a: banana}, {a: cat}]\n",
		},
	},
	{
		description:    "Sort array of objects by key",
		subdescription: "Note that you can give sort_by complex expressions, not just paths",
		document:       "cool: [{b: banana},{a: banana},{c: banana}]",
		expression:     `.cool |= sort_by(keys | .[0])`,
		expected: []string{
			"D0, P[], (doc)::cool: [{a: banana}, {b: banana}, {c: banana}]\n",
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
}

func TestSortByOperatorScenarios(t *testing.T) {
	for _, tt := range sortByOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "sort", sortByOperatorScenarios)
}
