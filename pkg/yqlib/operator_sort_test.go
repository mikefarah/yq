package yqlib

import "testing"

var sortByOperatorScenarios = []expressionScenario{
	{
		description: "Sort by string field",
		document:    "[{a: banana},{a: cat},{a: apple}]",
		expression:  `sort_by(.a)`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	// {
	// 	description: "Sort, nulls come first",
	// 	document:    "[8,3,null,6]",
	// 	expression:  `sort`,
	// 	expected: []string{
	// 		"D0, P[], (!!bool)::[null,3,6,8]\n",
	// 	},
	// },
}

func TestSortByOperatorScenarios(t *testing.T) {
	for _, tt := range sortByOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Sort", sortByOperatorScenarios)
}
