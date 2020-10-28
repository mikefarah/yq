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
	{
		document:   `{name: Mike, pets: [cat, dog], food: [hotdog, burger]}`,
		expression: `.name: .pets[], "f":.food[]`,
		expected: []string{
			"D0, P[], (!!seq)::- Mike: cat\n- Mike: dog\n",
			"D0, P[], (!!seq)::- f: hotdog\n- f: burger\n",
		},
	},
	{
		document:   `{name: Mike, pets: {cows: [apl, bba]}}`,
		expression: `"a":.name, "b":.pets`,
		expected: []string{
			"D0, P[], (!!seq)::- a: Mike\n",
			"D0, P[], (!!seq)::- b: {cows: [apl, bba]}\n",
		},
	},
}

func TestCreateMapOperatorScenarios(t *testing.T) {
	for _, tt := range createMapOperatorScenarios {
		testScenario(t, &tt)
	}
}
