package yqlib

import (
	"testing"
)

var deleteOperatorScenarios = []expressionScenario{
	{
		description: "Delete entry in map",
		document:    `{a: cat, b: dog}`,
		expression:  `del(.b)`,
		expected: []string{
			"D0, P[], (doc)::{a: cat}\n",
		},
	},
	{
		description: "Delete nested entry in map",
		document:    `{a: {a1: fred, a2: frood}}`,
		expression:  `del(.a.a1)`,
		expected: []string{
			"D0, P[], (doc)::{a: {a2: frood}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {a1: fred, a2: frood}}`,
		expression: `.a | del(.a1)`,
		expected: []string{
			"D0, P[a], (!!map)::{a2: frood}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: [1,2,3]`,
		expression: `.a | del(.[1])`,
		expected: []string{
			"D0, P[a], (!!seq)::[1, 3]\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[0, {a: cat, b: dog}]`,
		expression: `.[1] | del(.a)`,
		expected: []string{
			"D0, P[1], (!!map)::{b: dog}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[{a: cat, b: dog}]`,
		expression: `.[0] | del(.a)`,
		expected: []string{
			"D0, P[0], (!!map)::{b: dog}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[{a: {b: thing, c: frog}}]`,
		expression: `.[0].a | del(.b)`,
		expected: []string{
			"D0, P[0 a], (!!map)::{c: frog}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[{a: {b: thing, c: frog}}]`,
		expression: `.[0] | del(.a.b)`,
		expected: []string{
			"D0, P[0], (!!map)::{a: {c: frog}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: [0, {b: thing, c: frog}]}`,
		expression: `.a[1] | del(.b)`,
		expected: []string{
			"D0, P[a 1], (!!map)::{c: frog}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: [0, {b: thing, c: frog}]}`,
		expression: `.a | del(.[1].b)`,
		expected: []string{
			"D0, P[a], (!!seq)::[0, {c: frog}]\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {a1: fred, a2: frood}}`,
		expression: `del(.. | select(.=="frood"))`,
		expected: []string{
			"D0, P[], (!!map)::{a: {a1: fred}}\n",
		},
	},
	{
		description: "Delete entry in array",
		document:    `[1,2,3]`,
		expression:  `del(.[1])`,
		expected: []string{
			"D0, P[], (doc)::[1, 3]\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: [1,2,3]`,
		expression: `del(.a[])`,
		expected: []string{
			"D0, P[], (doc)::a: []\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: [10,x,10, 10, x, 10]`,
		expression: `del(.a[] | select(. == 10))`,
		expected: []string{
			"D0, P[], (doc)::a: [x, x]\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: {thing1: yep, thing2: cool, thing3: hi, b: {thing1: cool, great: huh}}`,
		expression: `del(..)`,
		expected: []string{
			"D0, P[], (!!map)::{}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: {thing1: yep, thing2: cool, thing3: hi, b: {thing1: cool, great: huh}}`,
		expression: `del(.. | select(tag == "!!map") | (.b.thing1,.thing2))`,
		expected: []string{
			"D0, P[], (!!map)::a: {thing1: yep, thing3: hi, b: {great: huh}}\n",
		},
	},
	{
		description: "Delete nested entry in array",
		document:    `[{a: cat, b: dog}]`,
		expression:  `del(.[0].a)`,
		expected: []string{
			"D0, P[], (doc)::[{b: dog}]\n",
		},
	},
	{
		description: "Delete no matches",
		document:    `{a: cat, b: dog}`,
		expression:  `del(.c)`,
		expected: []string{
			"D0, P[], (doc)::{a: cat, b: dog}\n",
		},
	},
	{
		description: "Delete matching entries",
		document:    `{a: cat, b: dog, c: bat}`,
		expression:  `del( .[] | select(. == "*at") )`,
		expected: []string{
			"D0, P[], (doc)::{b: dog}\n",
		},
	},
	{
		description: "Recursively delete matching keys",
		document:    `{a: {name: frog, b: {name: blog, age: 12}}}`,
		expression:  `del(.. | select(has("name")).name)`,
		expected: []string{
			"D0, P[], (!!map)::{a: {b: {age: 12}}}\n",
		},
	},
}

func TestDeleteOperatorScenarios(t *testing.T) {
	for _, tt := range deleteOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "delete", deleteOperatorScenarios)
}
