package treeops

import (
	"testing"
)

var createMapOperatorScenarios = []expressionScenario{
	{
		document:   `{name: Mike, age: 32}`,
		expression: `.name: .age`,
		expected: []string{
			"D0, P[], (!!seq)::- Mike: 32\n",
		},
	},
	{
		document:   `{name: Mike, pets: [cat, dog]}`,
		expression: `.name: .pets[]`,
		expected: []string{
			"D0, P[], (!!seq)::- Mike: cat\n- Mike: dog\n",
		},
	},
}

func TestCreateMapOperatorScenarios(t *testing.T) {
	for _, tt := range createMapOperatorScenarios {
		testScenario(t, &tt)
	}
}
