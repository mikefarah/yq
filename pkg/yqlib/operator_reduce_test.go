package yqlib

import (
	"testing"
)

var reduceOperatorScenarios = []expressionScenario{
	{
		document:   `[10,2, 5, 3]`,
		expression: `.[] as $item reduce (0; . + $item)`,
		expected: []string{
			"D0, P[], (!!int)::20\n",
		},
	},
}

func TestReduceOperatorScenarios(t *testing.T) {
	for _, tt := range reduceOperatorScenarios {
		testScenario(t, &tt)
	}
	// documentScenarios(t, "Reduce", reduceOperatorScenarios)
}
