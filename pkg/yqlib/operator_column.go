package yqlib

import (
	"container/list"
	"fmt"
)

func columnOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("columnOperator")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		result := candidate.CreateReplacement()
		result.Kind = ScalarNode
		result.Value = fmt.Sprintf("%v", candidate.Column)
		result.Tag = "!!int"

		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}
