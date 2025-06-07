package yqlib

import (
	"testing"
)

var uniqueOperatorScenarios = []expressionScenario{
	{
		description:    "Unique array of scalars (string/numbers)",
		subdescription: "Note that unique maintains the original order of the array.",
		document:       `[2,1,3,2]`,
		expression:     `unique`,
		expected: []string{
			"D0, P[], (!!seq)::[2, 1, 3]\n",
		},
	},
	{
		description: "Unique splat",
		skipDoc:     true,
		document:    `[2,1,2]`,
		expression:  `unique[]`,
		expected: []string{
			"D0, P[0], (!!int)::2\n",
			"D0, P[1], (!!int)::1\n",
		},
	},
	{
		description:    "Unique nulls",
		subdescription: "Unique works on the node value, so it considers different representations of nulls to be different",
		document:       `[~,null, ~, null]`,
		expression:     `unique`,
		expected: []string{
			"D0, P[], (!!seq)::[~, null]\n",
		},
	},
	{
		description:    "Unique all nulls",
		subdescription: "Run against the node tag to unique all the nulls",
		document:       `[~,null, ~, null]`,
		expression:     `unique_by(tag)`,
		expected: []string{
			"D0, P[], (!!seq)::[~]\n",
		},
	},
	{
		description: "Unique array objects",
		document:    `[{name: harry, pet: cat}, {name: billy, pet: dog}, {name: harry, pet: cat}]`,
		expression:  `unique`,
		expected: []string{
			"D0, P[], (!!seq)::[{name: harry, pet: cat}, {name: billy, pet: dog}]\n",
		},
	},
	{
		description: "Unique array of objects by a field",
		document:    `[{name: harry, pet: cat}, {name: billy, pet: dog}, {name: harry, pet: dog}]`,
		expression:  `unique_by(.name)`,
		expected: []string{
			"D0, P[], (!!seq)::[{name: harry, pet: cat}, {name: billy, pet: dog}]\n",
		},
	},
	{
		description: "Unique array of arrays",
		document:    `[[cat,dog], [cat, sheep], [cat,dog]]`,
		expression:  `unique`,
		expected: []string{
			"D0, P[], (!!seq)::[[cat, dog], [cat, sheep]]\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[{name: harry, pet: cat}, {pet: fish}, {name: harry, pet: dog}]`,
		expression: `unique_by(.name)`,
		expected: []string{
			"D0, P[], (!!seq)::[{name: harry, pet: cat}, {pet: fish}]\n",
		},
	},
	{
		description: "unique by splat",
		skipDoc:     true,
		document:    `[{name: harry, pet: cat}, {pet: fish}, {name: harry, pet: dog}]`,
		expression:  `unique_by(.name)[]`,
		expected: []string{
			"D0, P[0], (!!map)::{name: harry, pet: cat}\n",
			"D0, P[1], (!!map)::{pet: fish}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[{name: harry, pet: cat}, {pet: fish}, {name: harry, pet: dog}]`,
		expression: `unique_by(.cat.dog)`,
		expected: []string{
			"D0, P[], (!!seq)::[{name: harry, pet: cat}]\n",
		},
	},
	{
		skipDoc:    true,
		document:   "# abc\n[{name: harry, pet: cat}, {pet: fish}, {name: harry, pet: dog}]\n# xyz",
		expression: `unique_by(.name)`,
		expected: []string{
			"D0, P[], (!!seq)::# abc\n[{name: harry, pet: cat}, {pet: fish}]\n# xyz\n",
		},
		skipForGoccy: true, // https://github.com/goccy/go-yaml/issues/757
	},
}

func TestUniqueOperatorScenarios(t *testing.T) {
	for _, tt := range uniqueOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "unique", uniqueOperatorScenarios)
}
