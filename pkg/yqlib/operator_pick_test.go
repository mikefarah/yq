package yqlib

import (
	"testing"
)

var pickOperatorScenarios = []expressionScenario{
	{
		description:    "Pick keys from map",
		subdescription: "Note that the order of the keys matches the pick order and non existent keys are skipped.",
		document:       "myMap: {cat: meow, dog: bark, thing: hamster, hamster: squeak}\n",
		expression:     `.myMap |= pick(["hamster", "cat", "goat"])`,
		expected: []string{
			"D0, P[], (!!map)::myMap: {hamster: squeak, cat: meow}\n",
		},
	},
	{
		description:    "Pick keys from map, included all the keys",
		subdescription: "We create a map of the picked keys plus all the current keys, and run that through unique",
		document:       "myMap: {cat: meow, dog: bark, thing: hamster, hamster: squeak}\n",
		expression:     `.myMap |= pick( (["thing"] + keys) | unique)`,
		expected: []string{
			"D0, P[], (!!map)::myMap: {thing: hamster, cat: meow, dog: bark, hamster: squeak}\n",
		},
	},
	{
		description: "Pick splat",
		skipDoc:     true,
		document:    "{cat: meow, dog: bark, thing: hamster, hamster: squeak}\n",
		expression:  `pick(["hamster", "cat"])[]`,
		expected: []string{
			"D0, P[hamster], (!!str)::squeak\n",
			"D0, P[cat], (!!str)::meow\n",
		},
	},
	{
		description: "Pick keys from map",
		skipDoc:     true,
		document:    "!things myMap: {cat: meow, dog: bark, thing: hamster, hamster: squeak}\n",
		expression:  `.myMap |= pick(["hamster", "cat", "goat"])`,
		expected: []string{
			"D0, P[], (!!map)::!things myMap: {hamster: squeak, cat: meow}\n",
		},
	},
	{
		description: "Pick keys from map with comments",
		skipDoc:     true,
		document:    "# abc\nmyMap: {cat: meow, dog: bark, thing: hamster, hamster: squeak}\n# xyz\n",
		expression:  `.myMap |= pick(["hamster", "cat", "goat"])`,
		expected: []string{
			"D0, P[], (!!map)::# abc\nmyMap: {hamster: squeak, cat: meow}\n# xyz\n",
		},
	},
	{
		description:    "Pick indices from array",
		subdescription: "Note that the order of the indices matches the pick order and non existent indices are skipped.",
		document:       `[cat, leopard, lion]`,
		expression:     `pick([2, 0, 734, -5])`,
		expected: []string{
			"D0, P[], (!!seq)::[lion, cat]\n",
		},
	},
	{
		description: "Pick indices from array with comments",
		skipDoc:     true,
		document:    "# abc\n[cat, leopard, lion]\n# xyz",
		expression:  `pick([2, 0, 734, -5])`,
		expected: []string{
			"D0, P[], (!!seq)::# abc\n[lion, cat]\n# xyz\n",
		},
	},
}

func TestPickOperatorScenarios(t *testing.T) {
	for _, tt := range pickOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "pick", pickOperatorScenarios)
}
