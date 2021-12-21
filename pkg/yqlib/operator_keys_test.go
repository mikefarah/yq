package yqlib

import (
	"testing"
)

var keysOperatorScenarios = []expressionScenario{
	{
		description: "Map keys",
		document:    `{dog: woof, cat: meow}`,
		expression:  `keys`,
		expected: []string{
			"D0, P[], (!!seq)::- dog\n- cat\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{}`,
		expression: `keys`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n",
		},
	},
	{
		description: "Array keys",
		document:    `[apple, banana]`,
		expression:  `keys`,
		expected: []string{
			"D0, P[], (!!seq)::- 0\n- 1\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[]`,
		expression: `keys`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n",
		},
	},
	{
		description: "Retrieve array key",
		document:    "[1,2,3]",
		expression:  `.[1] | key`,
		expected: []string{
			"D0, P[1], (!!int)::1\n",
		},
	},
	{
		description: "Retrieve map key",
		document:    "a: thing",
		expression:  `.a | key`,
		expected: []string{
			"D0, P[a], (!!str)::a\n",
		},
	},
	{
		description: "No key",
		document:    "{}",
		expression:  `key`,
		expected:    []string{},
	},
	{
		description: "Update map key",
		document:    "a:\n  x: 3\n  y: 4",
		expression:  `(.a.x | key) = "meow"`,
		expected: []string{
			"D0, P[], (doc)::a:\n    meow: 3\n    y: 4\n",
		},
	},
	{
		description: "Get comment from map key",
		document:    "a: \n  # comment on key\n  x: 3\n  y: 4",
		expression:  `.a.x | key | headComment`,
		expected: []string{
			"D0, P[a x], (!!str)::comment on key\n",
		},
	},
}

func TestKeysOperatorScenarios(t *testing.T) {
	for _, tt := range keysOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "keys", keysOperatorScenarios)
}
