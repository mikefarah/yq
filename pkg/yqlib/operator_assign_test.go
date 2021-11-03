package yqlib

import (
	"testing"
)

var assignOperatorScenarios = []expressionScenario{
	{
		description: "Create yaml file",
		expression:  `.a.b = "cat" | .x = "frog"`,
		expected: []string{
			"D0, P[], ()::a:\n    b: cat\nx: frog\n",
		},
	},
	{
		skipDoc:    true,
		document:   "{}",
		expression: `.a |= .b`,
		expected: []string{
			"D0, P[], (doc)::{a: null}\n",
		},
	},
	{
		skipDoc:    true,
		document:   "{}",
		expression: `.a = .b`,
		expected: []string{
			"D0, P[], (doc)::{a: null}\n",
		},
	},
	{
		description: "Update node to be the child value",
		document:    `{a: {b: {g: foof}}}`,
		expression:  `.a |= .b`,
		expected: []string{
			"D0, P[], (doc)::{a: {g: foof}}\n",
		},
	},
	{
		description: "Double elements in an array",
		document:    `[1,2,3]`,
		expression:  `.[] |= . * 2`,
		expected: []string{
			"D0, P[], (doc)::[2, 4, 6]\n",
		},
	},
	{
		description:    "Update node from another file",
		subdescription: "Note this will also work when the second file is a scalar (string/number)",
		document:       `{a: apples}`,
		document2:      "{b: bob}",
		expression:     `select(fileIndex==0).a = select(fileIndex==1) | select(fileIndex==0)`,
		expected: []string{
			"D0, P[], (doc)::{a: {b: bob}}\n",
		},
	},
	{
		description: "Update node to be the sibling value",
		document:    `{a: {b: child}, b: sibling}`,
		expression:  `.a = .b`,
		expected: []string{
			"D0, P[], (doc)::{a: sibling, b: sibling}\n",
		},
	},
	{
		description: "Updated multiple paths",
		document:    `{a: fieldA, b: fieldB, c: fieldC}`,
		expression:  `(.a, .c) = "potatoe"`,
		expected: []string{
			"D0, P[], (doc)::{a: potatoe, b: fieldB, c: potatoe}\n",
		},
	},
	{
		description: "Update string value",
		document:    `{a: {b: apple}}`,
		expression:  `.a.b = "frog"`,
		expected: []string{
			"D0, P[], (doc)::{a: {b: frog}}\n",
		},
	},
	{
		description:    "Update string value via |=",
		subdescription: "Note there is no difference between `=` and `|=` when the RHS is a scalar",
		document:       `{a: {b: apple}}`,
		expression:     `.a.b |= "frog"`,
		expected: []string{
			"D0, P[], (doc)::{a: {b: frog}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {b: apple}}`,
		expression: `.a.b | (. |= "frog")`,
		expected: []string{
			"D0, P[a b], (!!str)::frog\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {b: apple}}`,
		expression: `.a.b |= 5`,
		expected: []string{
			"D0, P[], (doc)::{a: {b: 5}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {b: apple}}`,
		expression: `.a.b |= 3.142`,
		expected: []string{
			"D0, P[], (doc)::{a: {b: 3.142}}\n",
		},
	},
	{
		description:    "Update deeply selected results",
		subdescription: "Note that the LHS is wrapped in brackets! This is to ensure we dont first filter out the yaml and then update the snippet.",
		document:       `{a: {b: apple, c: cactus}}`,
		expression:     `(.a[] | select(. == "apple")) = "frog"`,
		expected: []string{
			"D0, P[], (doc)::{a: {b: frog, c: cactus}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {b: apple, c: cactus}}`,
		expression: `(.a.[] | select(. == "apple")) = "frog"`,
		expected: []string{
			"D0, P[], (doc)::{a: {b: frog, c: cactus}}\n",
		},
	},
	{
		description: "Update array values",
		document:    `[candy, apple, sandy]`,
		expression:  `(.[] | select(. == "*andy")) = "bogs"`,
		expected: []string{
			"D0, P[], (doc)::[bogs, apple, bogs]\n",
		},
	},
	{
		description:           "Update empty object",
		dontFormatInputForDoc: true,
		document:              `{}`,
		expression:            `.a.b |= "bogs"`,
		expected: []string{
			"D0, P[], (doc)::{a: {b: bogs}}\n",
		},
	},
	{
		description:           "Update empty object and array",
		dontFormatInputForDoc: true,
		document:              `{}`,
		expression:            `.a.b.[0] |= "bogs"`,
		expected: []string{
			"D0, P[], (doc)::{a: {b: [bogs]}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{}`,
		expression: `.a.b.[1].c |= "bogs"`,
		expected: []string{
			"D0, P[], (doc)::{a: {b: [null, {c: bogs}]}}\n",
		},
	},
}

func TestAssignOperatorScenarios(t *testing.T) {
	for _, tt := range assignOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "assign-update", assignOperatorScenarios)
}
