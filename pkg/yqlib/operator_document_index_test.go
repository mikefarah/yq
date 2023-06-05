package yqlib

import (
	"testing"
)

var documentIndexScenarios = []expressionScenario{
	{
		description: "Retrieve a document index",
		document:    "a: cat\n---\na: frog\n",
		expression:  `.a | document_index`,
		expected: []string{
			"D0, P[a], (!!int)::0\n",
			"D1, P[a], (!!int)::1\n",
		},
	},
	{
		description: "Retrieve a document index, shorthand",
		document:    "a: cat\n---\na: frog\n",
		expression:  `.a | di`,
		expected: []string{
			"D0, P[a], (!!int)::0\n",
			"D1, P[a], (!!int)::1\n",
		},
	},
	{
		description: "Filter by document index",
		document:    "a: cat\n---\na: frog\n",
		expression:  `select(document_index == 1)`,
		expected: []string{
			"D1, P[], (!!map)::a: frog\n",
		},
	},
	{
		description: "Filter by document index shorthand",
		document:    "a: cat\n---\na: frog\n",
		expression:  `select(di == 1)`,
		expected: []string{
			"D1, P[], (!!map)::a: frog\n",
		},
	},
	{
		description: "Print Document Index with matches",
		document:    "a: cat\n---\na: frog\n",
		expression:  `.a | ({"match": ., "doc": document_index})`,
		expected: []string{
			"D0, P[], (!!map)::match: cat\ndoc: 0\n",
			"D1, P[], (!!map)::match: frog\ndoc: 1\n",
		},
	},
}

func TestDocumentIndexScenarios(t *testing.T) {
	for _, tt := range documentIndexScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "document-index", documentIndexScenarios)
}
