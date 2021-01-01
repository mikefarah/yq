package yqlib

import (
	"testing"
)

var fileOperatorScenarios = []expressionScenario{
	{
		description: "Get filename",
		document:    `{a: cat}`,
		expression:  `filename`,
		expected: []string{
			"D0, P[], (!!str)::sample.yml\n",
		},
	},
	{
		description: "Get file index",
		document:    `{a: cat}`,
		expression:  `fileIndex`,
		expected: []string{
			"D0, P[], (!!int)::0\n",
		},
	},
	{
		description: "Get file index alias",
		document:    `{a: cat}`,
		expression:  `fi`,
		expected: []string{
			"D0, P[], (!!int)::0\n",
		},
	},
}

func TestFileOperatorsScenarios(t *testing.T) {
	for _, tt := range fileOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "File Operators", fileOperatorScenarios)
}
