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

		// iterate backwards so that array indices stay valid when deleting
		for childEl := splatted.MatchingNodes.Back(); childEl != nil; childEl = childEl.Prev() {
			child := childEl.Value.(*CandidateNode)

			rhs, err := d.GetMatchingNodes(splatted.SingleChildContext(child), expressionNode.RHS)
			if err != nil {
				return Context{}, err
			}

			first := rhs.MatchingNodes.Front()
			if first != nil {
				child.UpdateFrom(first.Value.(*CandidateNode), assignPreferences{})
				continue
			}

			// like jq, entries that the expression returns nothing for are removed
			switch candidate.Kind {
			case MappingNode:
				deleteFromMap(candidate, child.Key.Value)
			case SequenceNode:
				deleteFromArray(candidate, child.Key.Value)
			}
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
		if splatted.MatchingNodes.Len() == 0 {
			results.PushBack(candidate.Copy())
			continue
		}

		result, err := d.GetMatchingNodes(splatted, expressionNode.RHS)
		log.Debugf("expressionNode.Rhs %v", expressionNode.RHS.Operation.OperationType)
		log.Debugf("result %v", result)
		if err != nil {
			return Context{}, err
		}

		selfExpression := &ExpressionNode{Operation: &Operation{OperationType: selfReferenceOpType}}
		collected, err := collectTogether(d, result, selfExpression)
		if err != nil {
			return Context{}, err
		}
		collected.Style = candidate.Style

		results.PushBack(collected)

	}

	return context.ChildContext(results), nil
}
