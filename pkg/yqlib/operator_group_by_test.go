package yqlib

import (
	"testing"
)

var groupByOperatorScenarios = []expressionScenario{
	{
		description: "Group by field",
		document:    `[{foo: 1, bar: 10}, {foo: 3, bar: 100}, {foo: 1, bar: 1}]`,
		expression:  `group_by(.foo)`,
		expected: []string{
			"D0, P[], (!!seq)::- - {foo: 1, bar: 10}\n  - {foo: 1, bar: 1}\n- - {foo: 3, bar: 100}\n",
		},
	},
	{
		description: "Group by field, with nuls",
		document:    `[{cat: dog}, {foo: 1, bar: 10}, {foo: 3, bar: 100}, {no: foo for you}, {foo: 1, bar: 1}]`,
		expression:  `group_by(.foo)`,
		expected: []string{
			"D0, P[], (!!seq)::- - {cat: dog}\n  - {no: foo for you}\n- - {foo: 1, bar: 10}\n  - {foo: 1, bar: 1}\n- - {foo: 3, bar: 100}\n",
		},
	},
}

func TestGroupByOperatorScenarios(t *testing.T) {
	for _, tt := range groupByOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "group-by", groupByOperatorScenarios)
}
