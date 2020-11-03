package yqlib

import (
	"testing"
)

var notOperatorScenarios = []expressionScenario{
	// {
	// 	document:   `cat`,
	// 	expression: `. | not`,
	// 	expected: []string{
	// 		"D0, P[], (!!bool)::false\n",
	// 	},
	// },
	// {
	// 	document:   `1`,
	// 	expression: `. | not`,
	// 	expected: []string{
	// 		"D0, P[], (!!bool)::false\n",
	// 	},
	// },
	// {
	// 	document:   `0`,
	// 	expression: `. | not`,
	// 	expected: []string{
	// 		"D0, P[], (!!bool)::false\n",
	// 	},
	// },
	{
		document:   `~`,
		expression: `. | not`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	// {
	// 	document:   `false`,
	// 	expression: `. | not`,
	// 	expected: []string{
	// 		"D0, P[], (!!bool)::true\n",
	// 	},
	// },
	// {
	// 	document:   `true`,
	// 	expression: `. | not`,
	// 	expected: []string{
	// 		"D0, P[], (!!bool)::false\n",
	// 	},
	// },
}

func TestNotOperatorScenarios(t *testing.T) {
	for _, tt := range notOperatorScenarios {
		testScenario(t, &tt)
	}
}
