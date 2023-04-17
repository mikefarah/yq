package yqlib

import (
	"testing"
)

var columnOperatorScenarios = []expressionScenario{
	{
		description: "Returns column of _value_ node",
		document:    "a: cat\nb: bob",
		expression:  `.b | column`,
		expected: []string{
			"D0, P[b], (!!int)::4\n",
		},
	},
	{
		description:    "Returns column of _key_ node",
		subdescription: "Pipe through the key operator to get the column of the key",
		document:       "a: cat\nb: bob",
		expression:     `.b | key | column`,
		expected: []string{
			"D0, P[1], (!!int)::1\n",
		},
	},
	{
		description: "First column is 1",
		document:    "a: cat",
		expression:  `.a | key | column`,
		expected: []string{
			"D0, P[1], (!!int)::1\n",
		},
	},
	{
		description: "No column data is 0",
		expression:  `{"a": "new entry"} | column`,
		expected: []string{
			"D0, P[], (!!int)::0\n",
		},
	},
}

func TestColumnOperatorScenarios(t *testing.T) {
	for _, tt := range columnOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "column", columnOperatorScenarios)
}
