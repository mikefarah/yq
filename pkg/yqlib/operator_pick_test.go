package yqlib

import (
	"testing"
)

var pickOperatorScenarios = []expressionScenario{
	{
		description:    "Pick keys from map",
		subdescription: "Note that the order of the keys matches the pick order and non existent keys are skipped.",
		document:       "myMap: {cat: meow, dog: bark, thing: hamster, hamster: squeek}\n",
		expression:     `.myMap |= pick(["hamster", "cat", "goat"])`,
		expected: []string{
			"D0, P[], (doc)::myMap: {hamster: squeek, cat: meow}\n",
		},
	},
	{
		description:    "Pick indices from array",
		subdescription: "Note that the order of the indexes matches the pick order and non existent indexes are skipped.",
		document:       `[cat, leopard, lion]`,
		expression:     `pick([2, 0, 734, -5])`,
		expected: []string{
			"D0, P[], (!!seq)::[lion, cat]\n",
		},
	},
}

func TestPickOperatorScenarios(t *testing.T) {
	for _, tt := range pickOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "pick", pickOperatorScenarios)
}
