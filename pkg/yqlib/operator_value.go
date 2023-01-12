package yqlib

import "container/list"

func referenceOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	return context.SingleChildContext(expressionNode.Operation.CandidateNode), nil
}

func valueOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debug("value = %v", expressionNode.Operation.CandidateNode.Node.Value)
	if context.MatchingNodes.Len() == 0 {
		clone, err := expressionNode.Operation.CandidateNode.Copy()
		if err != nil {
			return Context{}, err
		}
		return context.SingleChildContext(clone), nil
	}

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		clone, err := expressionNode.Operation.CandidateNode.Copy()
		if err != nil {
			return Context{}, err
		}
		results.PushBack(clone)
	}

	return context.ChildContext(results), nil
}
