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
}

func TestDeleteOperatorScenarios(t *testing.T) {
	for _, tt := range deleteOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Delete", deleteOperatorScenarios)
}
