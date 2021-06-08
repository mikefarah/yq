package yqlib

import (
	"testing"
)

var entriesOperatorScenarios = []expressionScenario{
	{
		description: "to_entries Map",
		document:    `{a: 1, b: 2}`,
		expression:  `to_entries`,
		expected: []string{
			"D0, P[], (!!seq)::- key: a\n  value: 1\n- key: b\n  value: 2\n",
		},
	},
	{
		description: "to_entries Array",
		document:    `[a, b]`,
		expression:  `to_entries`,
		expected: []string{
			"D0, P[], (!!seq)::- key: 0\n  value: a\n- key: 1\n  value: b\n",
		},
	},
	{
		description: "to_entries null",
		document:    `null`,
		expression:  `to_entries`,
		expected:    []string{},
	},
	{
		description: "from_entries map",
		document:    `{a: 1, b: 2}`,
		expression:  `to_entries | from_entries`,
		expected: []string{
			"D0, P[], (!!map)::a: 1\nb: 2\n",
		},
	},
	{
		description:    "from_entries with numeric key indexes",
		subdescription: "from_entries always creates a map, even for numeric keys",
		document:       `[a,b]`,
		expression:     `to_entries | from_entries`,
		expected: []string{
			"D0, P[], (!!map)::0: a\n1: b\n",
		},
	},
	{
		description: "Use with_entries to update keys",
		document:    `{a: 1, b: 2}`,
		expression:  `with_entries(.key |= "KEY_" + .)`,
		expected: []string{
			"D0, P[], (!!map)::KEY_a: 1\nKEY_b: 2\n",
		},
	},
	{
		description: "Use with_entries to filter the map",
		document:    `{a: { b: bird }, c: { d: dog }}`,
		expression:  `with_entries(select(.value | has("b")))`,
		expected: []string{
			"D0, P[], (!!map)::a: {b: bird}\n",
		},
	},
}

func TestEntriesOperatorScenarios(t *testing.T) {
	for _, tt := range entriesOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Entries", entriesOperatorScenarios)
}
