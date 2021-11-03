package yqlib

import (
	"testing"
)

var sortKeysOperatorScenarios = []expressionScenario{
	{
		description: "Sort keys of map",
		document:    `{c: frog, a: blah, b: bing}`,
		expression:  `sortKeys(.)`,
		expected: []string{
			"D0, P[], (doc)::{a: blah, b: bing, c: frog}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{c: frog}`,
		expression: `sortKeys(.d)`,
		expected: []string{
			"D0, P[], (doc)::{c: frog}\n",
		},
	},
	{
		description:    "Sort keys recursively",
		subdescription: "Note the array elements are left unsorted, but maps inside arrays are sorted",
		document:       `{bParent: {c: dog, array: [3,1,2]}, aParent: {z: donkey, x: [{c: yum, b: delish}, {b: ew, a: apple}]}}`,
		expression:     `sortKeys(..)`,
		expected: []string{
			"D0, P[], (!!map)::{aParent: {x: [{b: delish, c: yum}, {a: apple, b: ew}], z: donkey}, bParent: {array: [3, 1, 2], c: dog}}\n",
		},
	},
}

func TestSortKeysOperatorScenarios(t *testing.T) {
	for _, tt := range sortKeysOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "sort-keys", sortKeysOperatorScenarios)
}
