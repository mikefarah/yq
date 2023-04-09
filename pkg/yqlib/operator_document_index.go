package yqlib

import (
	"container/list"
	"fmt"
)

func getDocumentIndexOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		scalar := candidate.CreateReplacement(ScalarNode, "!!int", fmt.Sprintf("%v", candidate.Document))
		results.PushBack(scalar)
	}
	return context.ChildContext(results), nil
}
