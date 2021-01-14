package yqlib

import (
	"testing"
)

var stringsOperatorScenarios = []expressionScenario{
	{
		description: "Join strings",
		document:    `[cat, meow, 1, null, true]`,
		expression:  `join("; ")`,
		expected: []string{
			"D0, P[], (!!str)::cat; meow; 1; ; true\n",
		},
	},
	{
		description: "Split strings",
		document:    `"cat; meow; 1; ; true"`,
		expression:  `split("; ")`,
		expected: []string{
			"D0, P[], (!!seq)::- cat\n- meow\n- \"1\"\n- \"\"\n- \"true\"\n",
		},
	},
	{
		description: "Split strings one match",
		document:    `"word"`,
		expression:  `split("; ")`,
		expected: []string{
			"D0, P[], (!!seq)::- word\n",
		},
	},
	{
		skipDoc:    true,
		document:   `""`,
		expression: `split("; ")`,
		expected: []string{
			"D0, P[], (!!seq)::[]\n", // dont actually want this, just not to error
		},
	},
	{
		skipDoc:    true,
		expression: `split("; ")`,
		expected:   []string{},
	},
}

func TestStringsOperatorScenarios(t *testing.T) {
	for _, tt := range stringsOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "String Operators", stringsOperatorScenarios)
}
