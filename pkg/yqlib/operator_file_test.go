package yqlib

import (
	"testing"
)

var fileOperatorScenarios = []expressionScenario{
	{
		description: "Get filename",
		document:    `{}`,
		expression:  `filename`,
		expected: []string{
			"D0, P[], (!!str)::sample.yml\n",
		},
	},
	{
		description: "Get file index",
		document:    `{}`,
		expression:  `fileIndex`,
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
