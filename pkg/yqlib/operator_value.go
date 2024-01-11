package yqlib

import "container/list"

func referenceOperator(_ *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	return context.SingleChildContext(expressionNode.Operation.CandidateNode), nil
}

func valueOperator(_ *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debug("value = %v", expressionNode.Operation.CandidateNode.Value)
	if context.MatchingNodes.Len() == 0 {
		clone := expressionNode.Operation.CandidateNode.Copy()
		return context.SingleChildContext(clone), nil
	}

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		clone := expressionNode.Operation.CandidateNode.Copy()
		results.PushBack(clone)
	}

	return context.ChildContext(results), nil
}
