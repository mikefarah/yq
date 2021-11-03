package yqlib

import (
	"testing"
)

var tagOperatorScenarios = []expressionScenario{
	{
		description: "Get tag",
		document:    `{a: cat, b: 5, c: 3.2, e: true, f: []}`,
		expression:  `.. | tag`,
		expected: []string{
			"D0, P[], (!!str)::'!!map'\n",
			"D0, P[a], (!!str)::'!!str'\n",
			"D0, P[b], (!!str)::'!!int'\n",
			"D0, P[c], (!!str)::'!!float'\n",
			"D0, P[e], (!!str)::'!!bool'\n",
			"D0, P[f], (!!str)::'!!seq'\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: cat, b: 5, c: 3.2, e: true, f: []}`,
		expression: `tag`,
		expected: []string{
			"D0, P[], (!!str)::'!!map'\n",
		},
	},
	{
		skipDoc:    true,
		document:   `32`,
		expression: `. tag= "!!str"`,
		expected: []string{
			"D0, P[], (doc)::\"32\"\n",
		},
	},
	{
		description: "Set custom tag",
		document:    `{a: str}`,
		expression:  `.a tag = "!!mikefarah"`,
		expected: []string{
			"D0, P[], (doc)::{a: !!mikefarah str}\n",
		},
	},
	{
		description: "Find numbers and convert them to strings",
		document:    `{a: cat, b: 5, c: 3.2, e: true}`,
		expression:  `(.. | select(tag == "!!int")) tag= "!!str"`,
		expected: []string{
			"D0, P[], (!!map)::{a: cat, b: \"5\", c: 3.2, e: true}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: "!!frog", b: "!!customTag"}`,
		expression: `.[] tag |= .`,
		expected: []string{
			"D0, P[], (doc)::{a: !!frog \"!!frog\", b: !!customTag \"!!customTag\"}\n",
		},
	},
}

func TestTagOperatorScenarios(t *testing.T) {
	for _, tt := range tagOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "tag", tagOperatorScenarios)
}
