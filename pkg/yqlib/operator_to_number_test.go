package yqlib

import (
	"testing"
)

var toNumberScenarios = []expressionScenario{
	{
		description: "Converts strings to numbers",
		document:    `["3", "3.1", "-1e3"]`,
		expression:  `.[] | to_number`,
		expected: []string{
			"D0, P[0], (!!int)::3\n",
			"D0, P[1], (!!float)::3.1\n",
			"D0, P[2], (!!float)::-1e3\n",
		},
	},
	{
		skipDoc:     true,
		description: "Converts strings to numbers, with tonumber because jq",
		document:    `["3", "3.1", "-1e3"]`,
		expression:  `.[] | tonumber`,
		expected: []string{
			"D0, P[0], (!!int)::3\n",
			"D0, P[1], (!!float)::3.1\n",
			"D0, P[2], (!!float)::-1e3\n",
		},
	},
	{
		description: "Doesn't change numbers",
		document:    `[3, 3.1, -1e3]`,
		expression:  `.[] | to_number`,
		expected: []string{
			"D0, P[0], (!!int)::3\n",
			"D0, P[1], (!!float)::3.1\n",
			"D0, P[2], (!!float)::-1e3\n",
		},
	},
	{
		skipDoc:     true,
		description: "Leading zero integers are decimal ints, not floats",
		document:    `[090, 018, -018, +012, 00, 08.5, !!float 090, "090"]`,
		expression:  `.[] | tag`,
		expected: []string{
			"D0, P[0], (!!str)::!!int\n",
			"D0, P[1], (!!str)::!!int\n",
			"D0, P[2], (!!str)::!!int\n",
			"D0, P[3], (!!str)::!!int\n",
			"D0, P[4], (!!str)::!!int\n",
			"D0, P[5], (!!str)::!!float\n",
			"D0, P[6], (!!str)::!!float\n",
			"D0, P[7], (!!str)::!!str\n",
		},
	},
	{
		skipDoc:     true,
		description: "Leading zero integers can be used in arithmetic",
		document:    `a: 090`,
		expression:  `.a + 1`,
		expected: []string{
			"D0, P[a], (!!int)::91\n",
		},
	},
	{
		description:   "Cannot convert null",
		expression:    `.a.b | to_number`,
		expectedError: "cannot convert node value [null] at path a.b of tag !!null to number",
	},
}

func TestToNumberOperatorScenarios(t *testing.T) {
	for _, tt := range toNumberScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "to_number", toNumberScenarios)
}
