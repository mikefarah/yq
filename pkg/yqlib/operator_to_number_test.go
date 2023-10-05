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
