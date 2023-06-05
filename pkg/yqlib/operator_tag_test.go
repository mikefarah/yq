package yqlib

import (
	"testing"
)

var tagOperatorScenarios = []expressionScenario{
	{
		description:    "tag of key is not a key",
		subdescription: "so it should have 'a' as the path",
		skipDoc:        true,
		document:       "a: frog\n",
		expression:     `.a | key | tag`,
		expected: []string{
			"D0, P[a], (!!str)::!!str\n",
		},
	},
	{
		description: "Get tag",
		document:    `{a: cat, b: 5, c: 3.2, e: true, f: []}`,
		expression:  `.. | tag`,
		expected: []string{
			"D0, P[], (!!str)::!!map\n",
			"D0, P[a], (!!str)::!!str\n",
			"D0, P[b], (!!str)::!!int\n",
			"D0, P[c], (!!str)::!!float\n",
			"D0, P[e], (!!str)::!!bool\n",
			"D0, P[f], (!!str)::!!seq\n",
		},
	},
	{
		description: "type is an alias for tag",
		document:    `{a: cat, b: 5, c: 3.2, e: true, f: []}`,
		expression:  `.. | type`,
		expected: []string{
			"D0, P[], (!!str)::!!map\n",
			"D0, P[a], (!!str)::!!str\n",
			"D0, P[b], (!!str)::!!int\n",
			"D0, P[c], (!!str)::!!float\n",
			"D0, P[e], (!!str)::!!bool\n",
			"D0, P[f], (!!str)::!!seq\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: cat, b: 5, c: 3.2, e: true, f: []}`,
		expression: `tag`,
		expected: []string{
			"D0, P[], (!!str)::!!map\n",
		},
	},
	{
		skipDoc:    true,
		document:   `32`,
		expression: `. tag= "!!str"`,
		expected: []string{
			"D0, P[], (!!str)::32\n",
		},
	},
	{
		description: "Set custom tag",
		document:    `{a: str}`,
		expression:  `.a tag = "!!mikefarah"`,
		expected: []string{
			"D0, P[], (!!map)::{a: !!mikefarah str}\n",
		},
	},
	{
		skipDoc:     true,
		description: "Set custom type",
		document:    `{a: str}`,
		expression:  `.a type = "!!mikefarah"`,
		expected: []string{
			"D0, P[], (!!map)::{a: !!mikefarah str}\n",
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
			"D0, P[], (!!map)::{a: !!frog \"!!frog\", b: !!customTag \"!!customTag\"}\n",
		},
	},
}

func TestTagOperatorScenarios(t *testing.T) {
	for _, tt := range tagOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "tag", tagOperatorScenarios)
}
