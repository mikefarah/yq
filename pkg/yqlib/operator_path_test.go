package yqlib

import (
	"testing"
)

var pathOperatorScenarios = []expressionScenario{
	{
		description: "Map path",
		document:    `{a: {b: cat}}`,
		expression:  `.a.b | path`,
		expected: []string{
			"D0, P[a b], (!!seq)::- a\n- b\n",
		},
	},
	{
		skipDoc: true,
		document: `a:
  b:
    c:
    - 0
    - 1
    - 2
    - 3`,
		expression: `.a.b.c.[]`,
		expected: []string{
			"D0, P[a b c 0], (!!int)::0\n",
			"D0, P[a b c 1], (!!int)::1\n",
			"D0, P[a b c 2], (!!int)::2\n",
			"D0, P[a b c 3], (!!int)::3\n",
		},
	},
	{
		description: "Get map key",
		document:    `{a: {b: cat}}`,
		expression:  `.a.b | path | .[-1]`,
		expected: []string{
			"D0, P[a b -1], (!!str)::b\n",
		},
	},
	{
		description: "Array path",
		document:    `{a: [cat, dog]}`,
		expression:  `.a.[] | select(. == "dog") | path`,
		expected: []string{
			"D0, P[a 1], (!!seq)::- a\n- 1\n",
		},
	},
	{
		description: "Get array index",
		document:    `{a: [cat, dog]}`,
		expression:  `.a.[] | select(. == "dog") | path | .[-1]`,
		expected: []string{
			"D0, P[a 1 -1], (!!int)::1\n",
		},
	},
	{
		description: "Print path and value",
		document:    `{a: [cat, dog, frog]}`,
		expression:  `.a[] | select(. == "*og") | [{"path":path, "value":.}]`,
		expected: []string{
			"D0, P[a 1], (!!seq)::- path:\n    - a\n    - 1\n  value: dog\n",
			"D0, P[a 2], (!!seq)::- path:\n    - a\n    - 2\n  value: frog\n",
		},
	},
	{
		description: "Set path",
		document:    `{a: {b: cat}}`,
		expression:  `setpath(["a", "b"]; "things")`,
		expected: []string{
			"D0, P[], (doc)::{a: {b: things}}\n",
		},
	},
	{
		description: "Set on empty document",
		expression:  `setpath(["a", "b"]; "things")`,
		expected: []string{
			"D0, P[], ()::a:\n    b: things\n",
		},
	},
	{
		description: "Set array path",
		document:    `a: [cat, frog]`,
		expression:  `setpath(["a", 0]; "things")`,
		expected: []string{
			"D0, P[], (doc)::a: [things, frog]\n",
		},
	},
	{
		description: "Set array path empty",
		expression:  `setpath(["a", 0]; "things")`,
		expected: []string{
			"D0, P[], ()::a:\n    - things\n",
		},
	},
}

func TestPathOperatorsScenarios(t *testing.T) {
	for _, tt := range pathOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "path", pathOperatorScenarios)
}
