package yqlib

import "container/list"

func unionOperator(d *dataTreeNavigator, matchingNodes *list.List, expressionNode *ExpressionNode) (*list.List, error) {
	lhs, err := d.GetMatchingNodes(matchingNodes, expressionNode.Lhs)
	if err != nil {
		return nil, err
	}
	rhs, err := d.GetMatchingNodes(matchingNodes, expressionNode.Rhs)
	if err != nil {
		return nil, err
	}
	for el := rhs.Front(); el != nil; el = el.Next() {
		node := el.Value.(*CandidateNode)
		lhs.PushBack(node)
	}
	return lhs, nil
}
