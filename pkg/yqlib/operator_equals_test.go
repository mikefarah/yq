package yqlib

import (
	"testing"
)

var equalsOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		expression: ".a == .b",
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: cat`,
		expression: ".a == .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: cat`,
		expression: ".b == .a",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "cat",
		document2:  "dog",
		expression: "select(fi==0) == select(fi==1)",
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "{}",
		expression: "(.a == .b) as $x",
		expected: []string{
			"D0, P[], (doc)::{}\n",
		},
	},
	{
		skipDoc:    true,
		document:   "{}",
		expression: ".a == .b",
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "{}",
		expression: "(.a != .b) as $x",
		expected: []string{
			"D0, P[], (doc)::{}\n",
		},
	},
	{
		skipDoc:    true,
		document:   "{}",
		expression: ".a != .b",
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "{a: {b: 10}}",
		expression: "select(.c != null)",
		expected:   []string{},
	},
	{
		skipDoc:    true,
		document:   "{a: {b: 10}}",
		expression: "select(.d == .c)",
		expected: []string{
			"D0, P[], (doc)::{a: {b: 10}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   "{a: {b: 10}}",
		expression: "select(null == .c)",
		expected: []string{
			"D0, P[], (doc)::{a: {b: 10}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   "{a: { b: {things: \"\"}, f: [1], g: [] }}",
		expression: ".. | select(. == \"\")",
		expected: []string{
			"D0, P[a b things], (!!str)::\"\"\n",
		},
	},
	{
		description: "Match string",
		document:    `[cat,goat,dog]`,
		expression:  `.[] | (. == "*at")`,
		expected: []string{
			"D0, P[0], (!!bool)::true\n",
			"D0, P[1], (!!bool)::true\n",
			"D0, P[2], (!!bool)::false\n",
		},
	},
	{
		description: "Don't match string",
		document:    `[cat,goat,dog]`,
		expression:  `.[] | (. != "*at")`,
		expected: []string{
			"D0, P[0], (!!bool)::false\n",
			"D0, P[1], (!!bool)::false\n",
			"D0, P[2], (!!bool)::true\n",
		},
	},
	{
		description: "Match number",
		document:    `[3, 4, 5]`,
		expression:  `.[] | (. == 4)`,
		expected: []string{
			"D0, P[0], (!!bool)::false\n",
			"D0, P[1], (!!bool)::true\n",
			"D0, P[2], (!!bool)::false\n",
		},
	},
	{
		description: "Dont match number",
		document:    `[3, 4, 5]`,
		expression:  `.[] | (. != 4)`,
		expected: []string{
			"D0, P[0], (!!bool)::true\n",
			"D0, P[1], (!!bool)::false\n",
			"D0, P[2], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: { cat: {b: apple, c: whatever}, pat: {b: banana} }`,
		expression: `.a | (.[].b == "apple")`,
		expected: []string{
			"D0, P[a cat b], (!!bool)::true\n",
			"D0, P[a pat b], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   ``,
		expression: `null == null`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		description: "Match nulls",
		document:    ``,
		expression:  `null == ~`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		description: "Non exisitant key doesn't equal a value",
		document:    "a: frog",
		expression:  `select(.b != "thing")`,
		expected: []string{
			"D0, P[], (doc)::a: frog\n",
		},
	},
	{
		description: "Two non existant keys are equal",
		document:    "a: frog",
		expression:  `select(.b == .c)`,
		expected: []string{
			"D0, P[], (doc)::a: frog\n",
		},
	},
}

func TestEqualOperatorScenarios(t *testing.T) {
	for _, tt := range equalsOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "equals", equalsOperatorScenarios)
}
