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
}

func TestKeysOperatorScenarios(t *testing.T) {
	for _, tt := range keysOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Keys", keysOperatorScenarios)
}
