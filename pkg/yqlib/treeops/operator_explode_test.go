package treeops

import (
	"testing"
)

var explodeTest = []expressionScenario{
	{
		document:   `{a: mike}`,
		expression: `explode(.a)`,
		expected: []string{
			"D0, P[], (doc)::{a: mike}\n",
		},
	},
	{
		document:   `{f : {a: &a cat, b: *a}}`,
		expression: `explode(.f)`,
		expected: []string{
			"D0, P[], (doc)::{f: {a: cat, b: cat}}\n",
		},
	},
	{
		document:   mergeDocSample,
		expression: `.foo* | explode(.)`,
		expected: []string{
			"D0, P[], (doc)::{f: {a: cat, b: cat}}\n",
		},
	},
}

func TestExplodeOperatorScenarios(t *testing.T) {
	for _, tt := range explodeTest {
		testScenario(t, &tt)
	}
}
