package yqlib

import "container/list"

func pipeOperator(d *dataTreeNavigator, matchingNodes *list.List, expressionNode *ExpressionNode) (*list.List, error) {
	lhs, err := d.GetMatchingNodes(matchingNodes, expressionNode.Lhs)
	if err != nil {
		return nil, err
	}
	return d.GetMatchingNodes(lhs, expressionNode.Rhs)
}
