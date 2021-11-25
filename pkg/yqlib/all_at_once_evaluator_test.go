package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var evaluateNodesScenario = []expressionScenario{
	{
		document:   `a: hello`,
		expression: `.a`,
		expected: []string{
			"D0, P[a], (!!str)::hello\n",
		},
	},
	{
		document:   `a: hello`,
		expression: `.`,
		expected: []string{
			"D0, P[], (doc)::a: hello\n",
		},
	},
	{
		document:   `- a: "yes"`,
		expression: `.[] | has("a")`,
		expected: []string{
			"D0, P[0], (!!bool)::true\n",
		},
	},
}

func TestAllAtOnceEvaluateNodes(t *testing.T) {
	evaluator := NewAllAtOnceEvaluator()
	for _, tt := range evaluateNodesScenario {
		node := test.ParseData(tt.document)
		list, _ := evaluator.EvaluateNodes(tt.expression, &node)
		test.AssertResultComplex(t, tt.expected, resultsToString(t, list))
	}
}
