package yqlib

import (
	"testing"
)

var styleOperatorScenarios = []expressionScenario{
	{
		description: "Set tagged style",
		document:    `{a: cat, b: 5, c: 3.2, e: true}`,
		expression:  `.. style="tagged"`,
		expected: []string{
			"D0, P[], (doc)::{a: 'cat'}\n",
		},
	},
	{
		description: "Set double quote style",
		document:    `{a: cat, b: 5, c: 3.2, e: true}`,
		expression:  `.. style="double"`,
		expected: []string{
			"D0, P[], (doc)::{a: 'cat'}\n",
		},
	},
	{
		description: "Set single quote style",
		document:    `{a: cat, b: 5, c: 3.2, e: true}`,
		expression:  `.. style="single"`,
		expected: []string{
			"D0, P[], (doc)::{a: 'cat'}\n",
		},
	},
	{
		description: "Set literal quote style",
		document:    `{a: cat, b: 5, c: 3.2, e: true}`,
		expression:  `.. style="literal"`,
		expected: []string{
			"D0, P[], (doc)::{a: 'cat'}\n",
		},
	},
	{
		description: "Set folded quote style",
		document:    `{a: cat, b: 5, c: 3.2, e: true}`,
		expression:  `.. style="folded"`,
		expected: []string{
			"D0, P[], (doc)::{a: 'cat'}\n",
		},
	},
	{
		description: "Set flow quote style",
		document:    `{a: cat, b: 5, c: 3.2, e: true}`,
		expression:  `.. style="flow"`,
		expected: []string{
			"D0, P[], (doc)::{a: 'cat'}\n",
		},
	},
	{
		description: "Set empty (default) quote style",
		document:    `{a: cat, b: 5, c: 3.2, e: true}`,
		expression:  `.. style=""`,
		expected: []string{
			"D0, P[], (doc)::{a: 'cat'}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: cat, b: double}`,
		expression: `.a style=.b`,
		expected: []string{
			"D0, P[], (doc)::{a: \"cat\", b: double}\n",
		},
	},
	{
		description: "Read style",
		document:    `{a: "cat", b: 'thing'}`,
		expression:  `.. | style`,
		expected: []string{
			"D0, P[], (!!str)::flow\n",
			"D0, P[a], (!!str)::double\n",
			"D0, P[b], (!!str)::single\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: cat`,
		expression: `.. | style`,
		expected: []string{
			"D0, P[], (!!str)::\"\"\n",
			"D0, P[a], (!!str)::\"\"\n",
		},
	},
}

func TestStyleOperatorScenarios(t *testing.T) {
	for _, tt := range styleOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Style Operator", styleOperatorScenarios)
}
