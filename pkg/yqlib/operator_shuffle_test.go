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
		description: "Shuffle array in place",
		document:    "cool: [1, 2, 3, 4, 5]",
		expression:  `.cool |= shuffle`,
		expected: []string{
			"D0, P[], (doc)::cool: [5, 2, 4, 1, 3]\n",
		},
	},
}

func TestShuffleByOperatorScenarios(t *testing.T) {
	for _, tt := range shuffleOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "shuffle", shuffleOperatorScenarios)
}
