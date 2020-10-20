package treeops

import (
	"testing"
)

var traversePathOperatorScenarios = []expressionScenario{
	{
		document:   `{a: {b: apple}}`,
		expression: `.a`,
		expected: []string{
			"D0, P[a], (!!map)::{b: apple}\n",
		},
	},
	{
		document:   `[{b: apple}, {c: banana}]`,
		expression: `.[]`,
		expected: []string{
			"D0, P[0], (!!map)::{b: apple}\n",
			"D0, P[1], (!!map)::{c: banana}\n",
		},
	},
	{
		document:   `{}`,
		expression: `.a.b`,
		expected: []string{
			"D0, P[a b], ()::null\n",
		},
	},
	{
		document:   `{}`,
		expression: `.[1].a`,
		expected: []string{
			"D0, P[1 a], ()::null\n",
		},
	},
	{
		document:   `{}`,
		expression: `.a.[1]`,
		expected: []string{
			"D0, P[a 1], ()::null\n",
		},
	},
	{
		document:   `{a: {cat: apple, mad: things}}`,
		expression: `.a."*a*"`,
		expected: []string{
			"D0, P[a cat], (!!str)::apple\n",
			"D0, P[a mad], (!!str)::things\n",
		},
	},
	{
		document:   `{a: {cat: {b: 3}, mad: {b: 4}, fad: {c: t}}}`,
		expression: `.a."*a*".b`,
		expected: []string{
			"D0, P[a cat b], (!!int)::3\n",
			"D0, P[a mad b], (!!int)::4\n",
			"D0, P[a fad b], ()::null\n",
		},
	},
	{
		document:   `{a: {cat: apple, mad: things}}`,
		expression: `.a | (.cat, .mad)`,
		expected: []string{
			"D0, P[a cat], (!!str)::apple\n",
			"D0, P[a mad], (!!str)::things\n",
		},
	},
}

func TestTraversePathOperatorScenarios(t *testing.T) {
	for _, tt := range traversePathOperatorScenarios {
		testScenario(t, &tt)
	}
}
