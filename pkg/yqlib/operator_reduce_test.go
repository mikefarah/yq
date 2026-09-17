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
		description: "Merge all yaml files together",
		document:    `a: cat`,
		document2:   `b: dog`,
		expression:  `. as $item ireduce ({}; . * $item )`,
		expected: []string{
			"D0, P[], (!!map)::a: cat\nb: dog\n",
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
		skipDoc:     true,
		description: "Issue 2859: match then ireduce .captures[] matches jq reduce / capture",
		expression:  `"\"a\"=\"b\" \"x\"=\"y\"" | [match("\"(?<key>.*?)\"=\"(?<value>.*?)\""; "g") | (.captures[] as $capture ireduce ({}; .[$capture.name] = $capture.string))]`,
		expected: []string{
			"D0, P[], (!!seq)::- key: a\n  value: b\n- key: x\n  value: y\n",
		},
	},
}

func TestReduceOperatorScenarios(t *testing.T) {
	for _, tt := range reduceOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "reduce", reduceOperatorScenarios)
}
