package yqlib

import (
	"testing"
)

var omitOperatorScenarios = []expressionScenario{
	{
		description:    "Omit keys from map",
		subdescription: "Note that non existent keys are skipped.",
		document:       "myMap: {cat: meow, dog: bark, thing: hamster, hamster: squeak}\n",
		expression:     `.myMap |= omit(["hamster", "cat", "goat"])`,
		expected: []string{
			"D0, P[], (!!map)::myMap: {dog: bark, thing: hamster}\n",
		},
	},
	{
		description: "Omit splat",
		skipDoc:     true,
		document:    "{cat: meow, dog: bark, hamster: squeak}\n",
		expression:  `omit(["dog"])[]`,
		expected: []string{
			"D0, P[cat], (!!str)::meow\n",
			"D0, P[hamster], (!!str)::squeak\n",
		},
	},
	{
		description: "Omit keys from map",
		skipDoc:     true,
		document:    "!things myMap: {cat: meow, dog: bark, thing: hamster, hamster: squeak}\n",
		expression:  `.myMap |= omit(["hamster", "cat", "goat"])`,
		expected: []string{
			"D0, P[], (!!map)::!things myMap: {dog: bark, thing: hamster}\n",
		},
	},
	{
		description: "Omit keys from map with comments",
		skipDoc:     true,
		document:    "# abc\nmyMap: {cat: meow, dog: bark, thing: hamster, hamster: squeak}\n# xyz\n",
		expression:  `.myMap |= omit(["hamster", "cat", "goat"])`,
		expected: []string{
			"D0, P[], (!!map)::# abc\nmyMap: {dog: bark, thing: hamster}\n# xyz\n",
		},
		skipForGoccy: true,
	},
	{
		description:    "Omit indices from array",
		subdescription: "Note that non existent indices are skipped.",
		document:       `[cat, leopard, lion]`,
		expression:     `omit([2, 0, 734, -5])`,
		expected: []string{
			"D0, P[], (!!seq)::[leopard]\n",
		},
	},
	{
		description: "Omit indices from array with comments",
		skipDoc:     true,
		document:    "# abc\n[cat, leopard, lion]\n# xyz",
		expression:  `omit([2, 0, 734, -5])`,
		expected: []string{
			"D0, P[], (!!seq)::# abc\n[leopard]\n# xyz\n",
		},
		skipForGoccy: true,
	},
}

func TestOmitOperatorScenarios(t *testing.T) {
	for _, tt := range omitOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "omit", omitOperatorScenarios)
}
