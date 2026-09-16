package yqlib

import (
	"testing"
)

var inOperatorScenarios = []expressionScenario{
	{
		description: "Check key exists in map using variable binding",
		document:    "a: 1\nb: 2\nc: 3\n",
		expression:  `. as $m | "a" | in($m)`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		description: "Check key does not exist in map",
		document:    "a: 1\nb: 2\nc: 3\n",
		expression:  `. as $m | "d" | in($m)`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		description: "Check value exists in array",
		document:    "- Tool\n- Food\n- Flower\n",
		expression:  `. as $m | "Food" | in($m)`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		description: "Check value does not exist in array",
		document:    "- Tool\n- Food\n- Flower\n",
		expression:  `. as $m | "Animal" | in($m)`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		description:    "Check in with select on array elements",
		subdescription: "Filter items whose type is in the given list",
		document:       "- {item: Pizza, type: Food}\n- {item: Rose, type: Flower}\n- {item: Hammer, type: Tool}\n",
		expression:     `.[] | select(.type | in(["Tool", "Food"]))`,
		expected: []string{
			"D0, P[0], (!!map)::{item: Pizza, type: Food}\n",
			"D0, P[2], (!!map)::{item: Hammer, type: Tool}\n",
		},
	},
	{
		description: "In with variable binding - found",
		document:    "a: 1\nb: 2\nc: 3\n",
		expression:  `. as $m | "b" | in($m)`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		description: "In with variable binding - not found",
		document:    "a: 1\nb: 2\nc: 3\n",
		expression:  `. as $m | "z" | in($m)`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "- one\n- two\n- three\n",
		expression: `. as $m | "one" | in($m)`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "- one\n- two\n- three\n",
		expression: `. as $m | "five" | in($m)`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "key: value\nother: stuff\n",
		expression: `. as $m | "key" | in($m)`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "key: value\nother: stuff\n",
		expression: `. as $m | "missing" | in($m)`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
}

func TestInOperatorScenarios(t *testing.T) {
	for _, tt := range inOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "in", inOperatorScenarios)
}
