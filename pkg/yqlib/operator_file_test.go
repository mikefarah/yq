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
		description: "Get file indices of multiple documents",
		document:    `{a: cat}`,
		document2:   `{a: cat}`,
		expression:  `fileIndex`,
		expected: []string{
			"D0, P[], (!!int)::0\n",
			"D0, P[], (!!int)::1\n",
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
	{
		skipDoc:    true,
		document:   "a: cat\nb: dog",
		expression: `.. lineComment |= filename`,
		expected: []string{
			"D0, P[], (!!map)::a: cat # sample.yml\nb: dog # sample.yml\n",
		},
	},
}

func TestFileOperatorsScenarios(t *testing.T) {
	for _, tt := range fileOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "file-operators", fileOperatorScenarios)
}
