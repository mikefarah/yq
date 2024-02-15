package yqlib

import (
	"container/list"
)

func filterOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("filterOperation")
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		children := context.SingleChildContext(candidate)
		splatted, err := splat(children, traversePreferences{})
		if err != nil {
			return Context{}, err
		}
		filtered, err := selectOperator(d, splatted, expressionNode)
		if err != nil {
			return Context{}, err
		}

		selfExpression := &ExpressionNode{Operation: &Operation{OperationType: selfReferenceOpType}}
		collected, err := collectTogether(d, filtered, selfExpression)
		if err != nil {
			return Context{}, err
		}
		collected.Style = candidate.Style
		results.PushBack(collected)
	}
	return context.ChildContext(results), nil
}
