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
	{
		description:    "Merge all documents together - using context",
		subdescription: "The _$context_ variable set by reduce lets you access the data outside the reduce block.",
		document:       `a: cat`,
		document2:      `b: dog`,
		expression:     `fi as $item ireduce ({}; . * ($context | select(fileIndex==$item)) )`,
		expected: []string{
			"D0, P[], (!!map)::a: cat\nb: dog\n",
		},
	},
	{
		description:    "Merge all documents together - without using context",
		subdescription: "`$context` is just a convenient variable that `reduce` sets, you can use your own for more control",
		document:       `c: {a: cat}`,
		document2:      `c: {b: dog}`,
		expression:     `.c as $root | fileIndex as $item ireduce ({}; . * ($root | select(fileIndex==$item)) )`,
		expected: []string{
			"D0, P[], (!!map)::{a: cat, b: dog}\n",
		},
	},
}

func TestReduceOperatorScenarios(t *testing.T) {
	for _, tt := range reduceOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Reduce", reduceOperatorScenarios)
}
