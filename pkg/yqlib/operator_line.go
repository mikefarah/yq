package yqlib

import (
	"container/list"
	"fmt"
)

func lineOperator(_ *dataTreeNavigator, context Context, _ *ExpressionNode) (Context, error) {
	log.Debugf("lineOperator")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		result := candidate.CreateReplacement(ScalarNode, "!!int", fmt.Sprintf("%v", candidate.Line))
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}
