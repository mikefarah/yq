package yqlib

import "testing"

var reverseOperatorScenarios = []expressionScenario{
	{
		description: "Reverse",
		document:    "[1, 2, 3]",
		expression:  `reverse`,
		expected: []string{
			"D0, P[], (!!seq)::[3, 2, 1]\n",
		},
	},
	{
		skipDoc:    true,
		document:   "[]",
		expression: `reverse`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n",
		},
	},
	{
		skipDoc:    true,
		document:   "[1]",
		expression: `reverse`,
		expected: []string{
			"D0, P[], (!!seq)::[1]\n",
		},
	},
	{
		skipDoc:    true,
		document:   "[1,2]",
		expression: `reverse`,
		expected: []string{
			"D0, P[], (!!seq)::[2, 1]\n",
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
}

func TestReverseOperatorScenarios(t *testing.T) {
	for _, tt := range reverseOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "reverse", reverseOperatorScenarios)
}
