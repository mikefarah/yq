package yqlib

import (
	"testing"
)

var uniqueOperatorScenarios = []expressionScenario{
	{
		description: "Unique array of scalars (string/numbers)",
		document:    `[1,2,3,2]`,
		expression:  `unique`,
		expected: []string{
			"D0, P[], (!!seq)::- 1\n- 2\n- 3\n",
		},
	},
	{
		description:    "Unique nulls",
		subdescription: "Unique works on the node value, so it considers different representations of nulls to be different",
		document:       `[~,null, ~, null]`,
		expression:     `unique`,
		expected: []string{
			"D0, P[], (!!seq)::- ~\n- null\n",
		},
	},
	{
		description:    "Unique all nulls",
		subdescription: "Run against the node tag to unique all the nulls",
		document:       `[~,null, ~, null]`,
		expression:     `unique_by(tag)`,
		expected: []string{
			"D0, P[], (!!seq)::- ~\n",
		},
	},
	{
		description: "Unique array object fields",
		document:    `[{name: harry, pet: cat}, {name: billy, pet: dog}, {name: harry, pet: dog}]`,
		expression:  `unique_by(.name)`,
		expected: []string{
			"D0, P[], (!!seq)::- {name: harry, pet: cat}\n- {name: billy, pet: dog}\n",
		},
	},
}

func TestUniqueOperatorScenarios(t *testing.T) {
	for _, tt := range uniqueOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Unique", uniqueOperatorScenarios)
}
