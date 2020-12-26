package yqlib

import (
	"testing"
)

var traverseArrayOperatorScenarios = []expressionScenario{
	{
		document:   `[a,b,c]`,
		expression: `.[]`,
		expected: []string{
			"D0, P[0], (!!str)::a\n",
			"D0, P[1], (!!str)::b\n",
			"D0, P[2], (!!str)::c\n",
		},
	},
	{
		document:   `[a,b,c]`,
		expression: `[]`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n",
		},
	},
	{
		document:   `{a: [a,b,c]}`,
		expression: `.a[0]`,
		expected: []string{
			"D0, P[a 0], (!!str)::a\n",
		},
	},
	{
		document:   `{a: [a,b,c]}`,
		expression: `.a[0, 2]`,
		expected: []string{
			"D0, P[a 0], (!!str)::a\n",
			"D0, P[a 2], (!!str)::c\n",
		},
	},
	{
		document:   `{a: [a,b,c]}`,
		expression: `.a.[0, 2]`,
		expected: []string{
			"D0, P[a 0], (!!str)::a\n",
			"D0, P[a 2], (!!str)::c\n",
		},
	},
	{
		document:   `{a: [a,b,c]}`,
		expression: `.a.[0]`,
		expected: []string{
			"D0, P[a 0], (!!str)::a\n",
		},
	},
	{
		document:   `{a: [a,b,c]}`,
		expression: `.a[-1]`,
		expected: []string{
			"D0, P[a -1], (!!str)::c\n",
		},
	},
	{
		document:   `{a: [a,b,c]}`,
		expression: `.a.[-1]`,
		expected: []string{
			"D0, P[a -1], (!!str)::c\n",
		},
	},
	{
		document:   `{a: [a,b,c]}`,
		expression: `.a[-2]`,
		expected: []string{
			"D0, P[a -2], (!!str)::b\n",
		},
	},
	{
		document:   `{a: [a,b,c]}`,
		expression: `.a.[-2]`,
		expected: []string{
			"D0, P[a -2], (!!str)::b\n",
		},
	},
	{
		document:   `{a: [a,b,c]}`,
		expression: `.a[]`,
		expected: []string{
			"D0, P[a 0], (!!str)::a\n",
			"D0, P[a 1], (!!str)::b\n",
			"D0, P[a 2], (!!str)::c\n",
		},
	},
	{
		document:   `{a: [a,b,c]}`,
		expression: `.a.[]`,
		expected: []string{
			"D0, P[a 0], (!!str)::a\n",
			"D0, P[a 1], (!!str)::b\n",
			"D0, P[a 2], (!!str)::c\n",
		},
	},
	{
		document:   `{a: [a,b,c]}`,
		expression: `.a | .[]`,
		expected: []string{
			"D0, P[a 0], (!!str)::a\n",
			"D0, P[a 1], (!!str)::b\n",
			"D0, P[a 2], (!!str)::c\n",
		},
	},
}

func TestTraverseArrayOperatorScenarios(t *testing.T) {
	for _, tt := range traverseArrayOperatorScenarios {
		testScenario(t, &tt)
	}
}
