package yqlib

import (
	"testing"
)

var reduceOperatorScenarios = []expressionScenario{
	{
		description: "Sum numbers",
		document:    `[10,2, 5, 3]`,
		expression:  `.[] as $item ireduce (0; . + $item)`,
		expected: []string{
			"D0, P[], (!!int)::20\n",
		},
	},
	{
		description: "Convert an array to an object",
		document:    `[{name: Cathy, has: apples},{name: Bob, has: bananas}]`,
		expression:  `.[] as $item ireduce ({}; .[$item | .name] = ($item | .has) )`,
		expected: []string{
			"D0, P[], (!!map)::Cathy: apples\nBob: bananas\n",
		},
	},
}

func TestReduceOperatorScenarios(t *testing.T) {
	for _, tt := range reduceOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Reduce", reduceOperatorScenarios)
}
