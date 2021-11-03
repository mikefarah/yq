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
		expression:  `.a.[] | select(. == "*og") | [{"path":path, "value":.}]`,
		expected: []string{`D0, P[], (!!seq)::- path:
    - a
    - 1
  value: dog
- path:
    - a
    - 2
  value: frog
`},
	},
}

func TestPathOperatorsScenarios(t *testing.T) {
	for _, tt := range pathOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "path", pathOperatorScenarios)
}
