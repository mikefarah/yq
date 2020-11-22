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
		description: "Array path",
		document:    `{a: [cat, dog]}`,
		expression:  `.a.[] | select(. == "dog") | path`,
		expected: []string{
			"D0, P[a 1], (!!seq)::- a\n- 1\n",
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
	documentScenarios(t, "Path Operator", pathOperatorScenarios)
}
