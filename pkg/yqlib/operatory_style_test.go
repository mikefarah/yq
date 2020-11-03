package yqlib

import (
	"testing"
)

var styleOperatorScenarios = []expressionScenario{
	{
		document:   `{a: cat}`,
		expression: `.a style="single"`,
		expected: []string{
			"D0, P[], (doc)::{a: 'cat'}\n",
		},
	},
	{
		document:   `{a: "cat", b: 'dog'}`,
		expression: `.. style=""`,
		expected: []string{
			"D0, P[], (!!map)::a: cat\nb: dog\n",
		},
	},
	{
		document:   `{a: "cat", b: 'thing'}`,
		expression: `.. | style`,
		expected: []string{
			"D0, P[], (!!str)::flow\n",
			"D0, P[a], (!!str)::double\n",
			"D0, P[b], (!!str)::single\n",
		},
	},
	{
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
}
