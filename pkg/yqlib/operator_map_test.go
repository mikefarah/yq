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
		description: "mapping against an empty array should do nothing",
		skipDoc:     true,
		document:    `[]`,
		document2:   `["cat"]`,
		expression:  `map(3)`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n",
			"D0, P[], (!!seq)::[3]\n",
		},
	},
	{
		description: "mapping against an empty array should do nothing",
		skipDoc:     true,
		document:    `[[], [5]]`,
		expression:  `.[] |= map(3)`,
		expected: []string{
			"D0, P[], (!!seq)::[[], [3]]\n",
		},
	},
	{
		description: "mapping against an empty array should do nothing #2",
		skipDoc:     true,
		document:    `[]`,
		document2:   `[5]`,
		expression:  `map(3 + .)`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n",
			"D0, P[], (!!seq)::[8]\n",
		},
	},
	{
		description: "mapping against an empty array should do nothing",
		skipDoc:     true,
		document:    `[[], [5]]`,
		expression:  `.[] |= map(3 + .)`,
		expected: []string{
			"D0, P[], (!!seq)::[[], [8]]\n",
		},
	},
	{
		skipDoc:    true,
		expression: `[] | map(. + 42)`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[1,2]`,
		expression: `map(. + 1)[]`,
		expected: []string{
			"D0, P[0], (!!int)::2\n",
			"D0, P[1], (!!int)::3\n",
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
		document:   `{}`,
		document2:  `{b: 12}`,
		expression: `map_values(3)`,
		expected: []string{
			"D0, P[], (!!map)::{}\n",
			"D0, P[], (!!map)::{b: 3}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{}`,
		document2:  `{b: 12}`,
		expression: `map_values(3 + .)`,
		expected: []string{
			"D0, P[], (!!map)::{}\n",
			"D0, P[], (!!map)::{b: 15}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: 1, b: 2, c: 3}`,
		document2:  `{x: 10, y: 20, z: 30}`,
		expression: `map_values(. + 1)`,
		expected: []string{
			"D0, P[], (!!map)::{a: 2, b: 3, c: 4}\n",
			"D0, P[], (!!map)::{x: 11, y: 21, z: 31}\n",
		},
	},
	{
		description: "map values splat",
		skipDoc:     true,
		document:    `{a: 1, b: 2}`,
		expression:  `map_values(. + 1)[]`,
		expected: []string{
			"D0, P[a], (!!int)::2\n",
			"D0, P[b], (!!int)::3\n",
		},
	},
	{
		description: "Map object values",
		document:    `{a: 1, b: 2, c: 3}`,
		expression:  `map_values(. + 1)`,
		expected: []string{
			"D0, P[], (!!map)::{a: 2, b: 3, c: 4}\n",
		},
	},
}

func TestMapOperatorScenarios(t *testing.T) {
	for _, tt := range mapOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "map", mapOperatorScenarios)
}
