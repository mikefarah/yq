package yqlib

import (
	"testing"
)

var arrayToMapScenarios = []expressionScenario{
	{
		description: "Simple example",
		document:    `cool: [null, null, hello]`,
		expression:  `.cool |= array_to_map`,
		expected: []string{
			"D0, P[], (doc)::cool:\n    2: hello\n",
		},
	},
}

func TestArrayToMapScenarios(t *testing.T) {
	for _, tt := range arrayToMapScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "array-to-map", arrayToMapScenarios)
}
