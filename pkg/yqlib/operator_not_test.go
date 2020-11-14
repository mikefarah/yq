package yqlib

import (
	"testing"
)

var notOperatorScenarios = []expressionScenario{
	{
		description: "Not true is false",
		expression:  `true | not`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		description: "Not false is true",
		expression:  `false | not`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		description: "String values considered to be true",
		expression:  `"cat" | not`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		description: "Empty string value considered to be true",
		expression:  `"" | not`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		description: "Numbers are considered to be true",
		expression:  `1 | not`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		description: "Zero is considered to be true",
		expression:  `0 | not`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{

		description: "Null is considered to be false",
		expression:  `~ | not`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
}

func TestNotOperatorScenarios(t *testing.T) {
	for _, tt := range notOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Not Operator", notOperatorScenarios)
}
