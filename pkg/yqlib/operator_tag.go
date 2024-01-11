package yqlib

import (
	"container/list"
)

func assignTagOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	log.Debugf("AssignTagOperator: %v")
	tag := ""

	if !expressionNode.Operation.UpdateAssign {
		rhs, err := d.GetMatchingNodes(context.ReadOnlyClone(), expressionNode.RHS)
		if err != nil {
			return Context{}, err
		}

		if rhs.MatchingNodes.Front() != nil {
			tag = rhs.MatchingNodes.Front().Value.(*CandidateNode).Value
		}
	}

	lhs, err := d.GetMatchingNodes(context, expressionNode.LHS)

	if err != nil {
		return Context{}, err
	}

	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		log.Debugf("Setting tag of : %v", candidate.GetKey())
		if expressionNode.Operation.UpdateAssign {
			rhs, err := d.GetMatchingNodes(context.SingleReadonlyChildContext(candidate), expressionNode.RHS)
			if err != nil {
				return Context{}, err
			}

			if rhs.MatchingNodes.Front() != nil {
				tag = rhs.MatchingNodes.Front().Value.(*CandidateNode).Value
			}
		}
		candidate.Tag = tag
	}

	return context, nil
}

func getTagOperator(_ *dataTreeNavigator, context Context, _ *ExpressionNode) (Context, error) {
	log.Debugf("GetTagOperator")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		result := candidate.CreateReplacement(ScalarNode, "!!str", candidate.Tag)
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}
