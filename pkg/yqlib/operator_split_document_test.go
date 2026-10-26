package yqlib

import (
	"testing"
)

var splitDocOperatorScenarios = []expressionScenario{
	{
		description: "Split after assigning a variable",
		document:    `[1, 2]`,
		expression:  `.[] | . as $_ | split_doc`,
		expected: []string{
			"D0, P[0], (!!int)::1\n",
			"D1, P[1], (!!int)::2\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[{a: cat}, {b: dog}]`,
		expression: `.[] as $item | $item | split_doc`,
		expected: []string{
			"D0, P[0], (!!map)::{a: cat}\n",
			"D1, P[1], (!!map)::{b: dog}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[[1, 2], [3, 4]]`,
		expression: `.[] as $row | $row[] as $item | $item | split_doc`,
		expected: []string{
			"D0, P[0], (!!int)::1\n",
			"D1, P[1], (!!int)::2\n",
			"D2, P[0], (!!int)::3\n",
			"D3, P[1], (!!int)::4\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[1, 2, 3]`,
		expression: `.[] as $item | $item | select(. != 2) | split_doc`,
		expected: []string{
			"D0, P[0], (!!int)::1\n",
			"D1, P[2], (!!int)::3\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[1, 2]`,
		expression: `.[] | . as $_ | split_doc | split_doc`,
		expected: []string{
			"D0, P[0], (!!int)::1\n",
			"D1, P[1], (!!int)::2\n",
		},
	},
	{
		description: "Split empty",
		document:    ``,
		expression:  `split_doc`,
		expected: []string{
			"D0, P[], (!!null)::\n",
		},
	},
	{
		description: "Split array",
		document:    `[{a: cat}, {b: dog}]`,
		expression:  `.[] | split_doc`,
		expected: []string{
			"D0, P[0], (!!map)::{a: cat}\n",
			"D1, P[1], (!!map)::{b: dog}\n",
		},
	},
	{
		description: "Split splat",
		skipDoc:     true,
		document:    `[{a: cat}, {b: dog}]`,
		expression:  `.[] | split_doc[]`,
		expected: []string{
			"D0, P[0 a], (!!str)::cat\n",
			"D1, P[1 b], (!!str)::dog\n",
		},
	},
}

func TestSplitDocOperatorScenarios(t *testing.T) {
	for _, tt := range splitDocOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "split-into-documents", splitDocOperatorScenarios)
}

func TestSplitDocIndependentEvaluations(t *testing.T) {
	navigator := NewDataTreeNavigator()
	expression, err := getExpressionParser().ParseExpression(`.[] | . as $_ | split_doc`)
	if err != nil {
		t.Fatal(err)
	}
	for range 2 {
		inputs, err := readDocument(`[1, 2]`, "sample.yml", 0)
		if err != nil {
			t.Fatal(err)
		}
		result, err := navigator.GetMatchingNodes(Context{MatchingNodes: inputs}, expression)
		if err != nil {
			t.Fatal(err)
		}
		var expectedIndex uint
		for el := result.MatchingNodes.Front(); el != nil; el = el.Next() {
			candidate := el.Value.(*CandidateNode)
			if candidate.GetDocument() != expectedIndex {
				t.Errorf("expected document %d, got %d", expectedIndex, candidate.GetDocument())
			}
			expectedIndex++
		}
		if expectedIndex != 2 {
			t.Errorf("expected two documents, got %d", expectedIndex)
		}
	}
}
