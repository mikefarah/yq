package yqlib

import "container/list"

func valueOperator(d *dataTreeNavigator, matchMap *list.List, expressionNode *ExpressionNode) (*list.List, error) {
	log.Debug("value = %v", expressionNode.Operation.CandidateNode.Node.Value)
	return nodeToMap(expressionNode.Operation.CandidateNode), nil
}
