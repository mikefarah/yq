package yqlib

import (
	"testing"
)

var sortKeysOperatorScenarios = []expressionScenario{
	{
		description: "Sort keys of map",
		document:    `{c: frog, a: blah, b: bing}`,
		expression:  `sort_keys(.)`,
		expected: []string{
			"D0, P[], (!!map)::{a: blah, b: bing, c: frog}\n",
		},
	},
	{
		description: "Sort keys of map",
		skipDoc:     true,
		document:    `{c: frog, a: zoo}`,
		expression:  `sort_keys(.)[]`,
		expected: []string{
			"D0, P[a], (!!str)::zoo\n",
			"D0, P[c], (!!str)::frog\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{c: frog}`,
		expression: `sort_keys(.d)`,
		expected: []string{
			"D0, P[], (!!map)::{c: frog}\n",
		},
	},
	{
		description:    "Sort keys recursively",
		subdescription: "Note the array elements are left unsorted, but maps inside arrays are sorted",
		document:       `{bParent: {c: dog, array: [3,1,2]}, aParent: {z: donkey, x: [{c: yum, b: delish}, {b: ew, a: apple}]}}`,
		expression:     `sort_keys(..)`,
		expected: []string{
			"D0, P[], (!!map)::{aParent: {x: [{b: delish, c: yum}, {a: apple, b: ew}], z: donkey}, bParent: {array: [3, 1, 2], c: dog}}\n",
		},
	},
}

func TestSortKeysOperatorScenarios(t *testing.T) {
	for _, tt := range sortKeysOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "sort-keys", sortKeysOperatorScenarios)
}
