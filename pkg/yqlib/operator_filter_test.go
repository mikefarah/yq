package yqlib

import (
	"testing"
)

var filterOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		document:   `[1,2,3]`,
		document2:  `[1, 2]`,
		expression: `filter(. < 3)`,
		expected: []string{
			"D0, P[], (!!seq)::[2, 3, 4]\n",
			"D0, P[], (!!seq)::[1, 2]\n",
		},
	},
	/*
	{
		description: "Map array",
		document:    `[1,2,3]`,
		expression:  `map(. + 1)`,
		expected: []string{
			"D0, P[], (!!seq)::[2, 3, 4]\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: 1, b: 2, c: 3}`,
		document2:  `{x: 10, y: 20, z: 30}`,
		expression: `map_values(. + 1)`,
		expected: []string{
			"D0, P[], (doc)::{a: 2, b: 3, c: 4}\n",
			"D0, P[], (doc)::{x: 11, y: 21, z: 31}\n",
		},
	},
	{
		description: "Map object values",
		document:    `{a: 1, b: 2, c: 3}`,
		expression:  `map_values(. + 1)`,
		expected: []string{
			"D0, P[], (doc)::{a: 2, b: 3, c: 4}\n",
		},
	},
	*/
}

func TestFilterOperatorScenarios(t *testing.T) {
	for _, tt := range filterOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "filter", filterOperatorScenarios)
}
