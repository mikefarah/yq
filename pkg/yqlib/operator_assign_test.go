package yqlib

import (
	"testing"
)

var assignOperatorScenarios = []expressionScenario{
	{
		description: "Update node to be the child value",
		document:    `{a: {b: {g: foof}}}`,
		expression:  `.a |= .b`,
		expected: []string{
			"D0, P[], (doc)::{a: {g: foof}}\n",
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
		expression:  `(.a, .c) |= "potatoe"`,
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
		description: "Update selected results",
		document:    `{a: {b: apple, c: cactus}}`,
		expression:  `.a.[] | select(. == "apple") |= "frog"`,
		expected: []string{
			"D0, P[], (doc)::{a: {b: frog, c: cactus}}\n",
		},
	},
	{
		description: "Update array values",
		document:    `[candy, apple, sandy]`,
		expression:  `.[] | select(. == "*andy") |= "bogs"`,
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
		expression:            `.a.b[0] |= "bogs"`,
		expected: []string{
			"D0, P[], (doc)::{a: {b: [bogs]}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{}`,
		expression: `.a.b[1].c |= "bogs"`,
		expected: []string{
			"D0, P[], (doc)::{a: {b: [null, {c: bogs}]}}\n",
		},
	},
}

func TestAssignOperatorScenarios(t *testing.T) {
	for _, tt := range assignOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Assign", assignOperatorScenarios)
}
