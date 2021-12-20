package yqlib

import (
	"container/list"
)

func mapValuesOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		//run expression against entries
		// splat toEntries and pipe it into Rhs
		splatted, err := splat(context.SingleChildContext(candidate), traversePreferences{})
		if err != nil {
			return Context{}, err
		}

		assignUpdateExp := &ExpressionNode{
			Operation: &Operation{OperationType: assignOpType, UpdateAssign: true},
			Rhs:       expressionNode.Rhs,
		}
		_, err = assignUpdateOperator(d, splatted, assignUpdateExp)
		if err != nil {
			return Context{}, err
		}

	}

	return context, nil
}

func mapOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		//run expression against entries
		// splat toEntries and pipe it into Rhs
		splatted, err := splat(context.SingleChildContext(candidate), traversePreferences{})
		if err != nil {
			return Context{}, err
		}

		result, err := d.GetMatchingNodes(splatted, expressionNode.Rhs)
		log.Debug("expressionNode.Rhs %v", expressionNode.Rhs.Operation.OperationType)
		log.Debug("result %v", result)
		if err != nil {
			return Context{}, err
		}

		selfExpression := &ExpressionNode{Operation: &Operation{OperationType: selfReferenceOpType}}
		collected, err := collectTogether(d, result, selfExpression)
		if err != nil {
			return Context{}, err
		}
		collected.Node.Style = unwrapDoc(candidate.Node).Style

		results.PushBack(collected)

	}

	return context.ChildContext(results), nil
}
