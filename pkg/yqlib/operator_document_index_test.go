package yqlib

import (
	"testing"
)

var documentIndexScenarios = []expressionScenario{
	{
		description: "Retrieve a document index",
		document:    "a: cat\n---\na: frog\n",
		expression:  `.a | documentIndex`,
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
		expression:  `select(documentIndex == 1)`,
		expected: []string{
			"D1, P[], (doc)::a: frog\n",
		},
	},
	{
		description: "Filter by document index shorthand",
		document:    "a: cat\n---\na: frog\n",
		expression:  `select(di == 1)`,
		expected: []string{
			"D1, P[], (doc)::a: frog\n",
		},
	},
	{
		description: "Print Document Index with matches",
		document:    "a: cat\n---\na: frog\n",
		expression:  `.a | ({"match": ., "doc": documentIndex})`,
		expected: []string{
			"D0, P[], (!!map)::match: cat\ndoc: 0\n",
			"D0, P[], (!!map)::match: frog\ndoc: 1\n",
		},
	},
}

func TestDocumentIndexScenarios(t *testing.T) {
	for _, tt := range documentIndexScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Document Index", documentIndexScenarios)
}
