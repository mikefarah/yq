package yqlib

func pipeOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	if expressionNode.LHS.Operation.OperationType == assignVariableOpType {
		return variableLoop(d, context, expressionNode)
	}

	lhs, err := d.GetMatchingNodes(context, expressionNode.LHS)
	if err != nil {
		return Context{}, err
	}
	rhs, err := d.GetMatchingNodes(lhs, expressionNode.RHS)
	if err != nil {
		return Context{}, err
	}
	return context.ChildContext(rhs.MatchingNodes), nil
}
