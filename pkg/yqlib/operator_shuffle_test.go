package yqlib

import "testing"

var shuffleOperatorScenarios = []expressionScenario{
	{
		description: "Shuffle array",
		document:    "[1, 2, 3, 4, 5]",
		expression:  `shuffle`,
		expected: []string{
			"D0, P[], (!!seq)::[5, 2, 4, 1, 3]\n",
		},
	},
	{
		description: "Shuffle array",
		skipDoc:     true,
		document:    "[1, 2, 3]",
		expression:  `shuffle[]`,
		expected: []string{
			"D0, P[0], (!!int)::3\n",
			"D0, P[1], (!!int)::1\n",
			"D0, P[2], (!!int)::2\n",
		},
	},

	{
		description: "Shuffle array in place",
		document:    "cool: [1, 2, 3, 4, 5]",
		expression:  `.cool |= shuffle`,
		expected: []string{
			"D0, P[], (!!map)::cool: [5, 2, 4, 1, 3]\n",
		},
	},
}

func TestShuffleByOperatorScenarios(t *testing.T) {
	for _, tt := range shuffleOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "shuffle", shuffleOperatorScenarios)
}
