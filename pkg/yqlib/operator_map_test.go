package yqlib

import (
	"testing"
)

var mapOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		document:   `[1,2,3]`,
		document2:  `[5,6,7]`,
		expression: `map(. + 1)`,
		expected: []string{
			"D0, P[], (!!seq)::[2, 3, 4]\n",
			"D0, P[], (!!seq)::[6, 7, 8]\n",
		},
	},
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
}

func TestMapOperatorScenarios(t *testing.T) {
	for _, tt := range mapOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "map", mapOperatorScenarios)
}
